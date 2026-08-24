// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	stderrors "errors"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/bborbe/errors"
	libsentry "github.com/bborbe/sentry"
	"github.com/bborbe/service"
	"github.com/golang/glog"

	s3utils "github.com/bborbe/s3-utils"
)

func main() {
	app := &application{}
	os.Exit(service.Main(context.Background(), app, &app.SentryDSN, &app.SentryProxy))
}

type application struct {
	Port           int           `required:"false" arg:"port"              env:"PORT"              usage:"port to listen"                                                     default:"9090"`
	InitialDelay   time.Duration `required:"false" arg:"initial-delay"     env:"INITIAL_DELAY"     usage:"initial time before processing starts"                              default:"1m"`
	SentryDSN      string        `required:"false" arg:"sentry-dsn"        env:"SENTRY_DSN"        usage:"Sentry DSN"                                                                        display:"length"`
	SentryProxy    string        `required:"false" arg:"sentry-proxy"      env:"SENTRY_PROXY"      usage:"Sentry Proxy"`
	SrcS3Url       string        `required:"true"  arg:"src-s3-url"        env:"SRC_S3_URL"        usage:"URL of source S3 server"`
	SrcS3AccessKey string        `required:"true"  arg:"src-s3-access-key" env:"SRC_S3_ACCESS_KEY" usage:"Access Key for source S3 server"                                                   display:"length"`
	SrcS3SecretKey string        `required:"true"  arg:"src-s3-secret-key" env:"SRC_S3_SECRET_KEY" usage:"Secret Key for source S3 server"                                                   display:"length"`
	SrcS3Region    string        `required:"false" arg:"src-s3-region"     env:"SRC_S3_REGION"     usage:"S3 region for SigV4 signing (empty = SDK default)"`
	SrcBucket      string        `required:"true"  arg:"src-bucket"        env:"SRC_BUCKET"        usage:"source bucket"`
	DstS3Url       string        `required:"true"  arg:"dst-s3-url"        env:"DST_S3_URL"        usage:"URL of destination S3 server"`
	DstS3AccessKey string        `required:"true"  arg:"dst-s3-access-key" env:"DST_S3_ACCESS_KEY" usage:"Access Key for destination S3 server"                                              display:"length"`
	DstS3SecretKey string        `required:"true"  arg:"dst-s3-secret-key" env:"DST_S3_SECRET_KEY" usage:"Secret Key for destination S3 server"                                              display:"length"`
	DstS3Region    string        `required:"false" arg:"dst-s3-region"     env:"DST_S3_REGION"     usage:"S3 region for SigV4 signing (empty = SDK default)"`
	DstBucket      string        `required:"true"  arg:"dst-bucket"        env:"DST_BUCKET"        usage:"destination bucket"`
	Concurrency    int           `required:"false" arg:"concurrency"       env:"CONCURRENCY"       usage:"number of objects to copy concurrently"                             default:"16"`
	SkipExisting   bool          `required:"false" arg:"skip-existing"     env:"SKIP_EXISTING"     usage:"skip objects already present in the destination with matching size"`
}

func (a *application) Run(
	ctx context.Context,
	sentryClient libsentry.Client,
) error {
	srcClient := s3utils.CreateS3Client(
		s3utils.URL(a.SrcS3Url),
		s3utils.AccessKey(a.SrcS3AccessKey),
		s3utils.SecretKey(a.SrcS3SecretKey),
		s3utils.WithRegion(a.SrcS3Region),
	)
	dstClient := s3utils.CreateS3Client(
		s3utils.URL(a.DstS3Url),
		s3utils.AccessKey(a.DstS3AccessKey),
		s3utils.SecretKey(a.DstS3SecretKey),
		s3utils.WithRegion(a.DstS3Region),
	)
	uploader := manager.NewUploader(dstClient, func(u *manager.Uploader) {
		u.PartSize = 64 * 1024 * 1024
		u.Concurrency = 8
	})

	glog.V(2).
		Infof("copy bucket %s to bucket %s started (concurrency %d)", a.SrcBucket, a.DstBucket, a.Concurrency)
	paginator := s3.NewListObjectsV2Paginator(srcClient, &s3.ListObjectsV2Input{
		Bucket: aws.String(a.SrcBucket),
	})
	sem := make(chan struct{}, a.Concurrency)
	var wg sync.WaitGroup
	var count int64
	var copyErr error
	var mu sync.Mutex
	getCopyErr := func() error {
		mu.Lock()
		defer mu.Unlock()
		return copyErr
	}
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return errors.Wrap(ctx, err, "list objects failed")
		}
		for _, object := range output.Contents {
			select {
			case <-ctx.Done():
				mu.Lock()
				if copyErr == nil {
					copyErr = ctx.Err()
				}
				mu.Unlock()
			default:
			}
			if getCopyErr() != nil {
				break
			}
			sem <- struct{}{}
			wg.Add(1)
			go func(key *string, size *int64) {
				defer wg.Done()
				defer func() { <-sem }()
				if err := copyObject(ctx, srcClient, dstClient, uploader, a.SrcBucket, a.DstBucket, key, size, a.SkipExisting); err != nil {
					mu.Lock()
					if copyErr == nil {
						copyErr = err
					}
					mu.Unlock()
					return
				}
				c := atomic.AddInt64(&count, 1)
				if c%100 == 0 {
					glog.V(2).Infof("copied %d objects", c)
				}
			}(object.Key, object.Size)
		}
		if getCopyErr() != nil {
			break
		}
	}
	wg.Wait()
	if copyErr != nil {
		return errors.Wrap(ctx, copyErr, "copy object failed")
	}
	glog.V(2).
		Infof("copy bucket %s to bucket %s completed: %d objects", a.SrcBucket, a.DstBucket, count)
	return nil
}

// isNotFoundError reports whether the error is a definitive "object absent"
// (S3 NotFound / NoSuchKey). Anything else (403, 500, network) is NOT absent.
func isNotFoundError(err error) bool {
	var notFound *types.NotFound
	if stderrors.As(err, &notFound) {
		return true
	}
	var apiErr smithy.APIError
	return stderrors.As(err, &apiErr) &&
		(apiErr.ErrorCode() == "NotFound" || apiErr.ErrorCode() == "NoSuchKey")
}

func copyObject(
	ctx context.Context,
	srcClient *s3.Client,
	dstClient *s3.Client,
	uploader *manager.Uploader,
	srcBucket string,
	dstBucket string,
	key *string,
	srcSize *int64,
	skipExisting bool,
) error {
	if skipExisting {
		head, err := dstClient.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(dstBucket),
			Key:    key,
		})
		if err != nil {
			// Only treat a definitive "not found" as absent; surface
			// permission/server errors (403, 500) instead of silently
			// re-copying, which would obscure the real problem.
			if !isNotFoundError(err) {
				return errors.Wrap(ctx, err, "head object failed")
			}
		} else if head.ContentLength != nil && srcSize != nil &&
			*head.ContentLength == *srcSize {
			glog.V(4).Infof("skip %s (already present, matching size)", *key)
			return nil
		}
	}
	var lastErr error
	for attempt := 1; attempt <= 5; attempt++ {
		lastErr = copyObjectOnce(ctx, srcClient, uploader, srcBucket, dstBucket, key)
		if lastErr == nil {
			return nil
		}
		if attempt > 1 {
			glog.V(2).Infof("copy %s failed (attempt %d/5): %v", *key, attempt, lastErr)
		}
		select {
		case <-ctx.Done():
			return errors.Wrap(ctx, ctx.Err(), "copy object cancelled")
		case <-time.After(time.Duration(attempt) * time.Second):
		}
	}
	return lastErr
}

func copyObjectOnce(
	ctx context.Context,
	srcClient *s3.Client,
	uploader *manager.Uploader,
	srcBucket string,
	dstBucket string,
	key *string,
) error {
	getOutput, err := srcClient.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(srcBucket),
		Key:    key,
	})
	if err != nil {
		return errors.Wrap(ctx, err, "get object failed")
	}
	defer func() {
		if err := getOutput.Body.Close(); err != nil {
			glog.V(4).Infof("close body of %s failed: %v", *key, err)
		}
	}()
	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(dstBucket),
		Key:    key,
		Body:   getOutput.Body,
	})
	if err != nil {
		return errors.Wrap(ctx, err, "upload object failed")
	}
	glog.V(4).Infof("copied %s", *key)
	return nil
}
