// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bborbe/errors"
	libsentry "github.com/bborbe/sentry"
	"github.com/bborbe/service"
	"github.com/golang/glog"

	"github.com/bborbe/s3-utils"
)

func main() {
	app := &application{}
	os.Exit(service.Main(context.Background(), app, &app.SentryDSN, &app.SentryProxy))
}

type application struct {
	Port         int           `required:"false" arg:"port" env:"PORT" usage:"port to listen" default:"9090"`
	InitialDelay time.Duration `required:"false" arg:"initial-delay" env:"INITIAL_DELAY" usage:"initial time before processing starts" default:"1m"`
	SentryDSN    string        `required:"false" arg:"sentry-dsn" env:"SENTRY_DSN" usage:"Sentry DSN"`
	SentryProxy  string        `required:"false" arg:"sentry-proxy" env:"SENTRY_PROXY" usage:"Sentry Proxy"`
	S3Url        string        `required:"false" arg:"s3-url" env:"S3_URL" usage:"URL of S3 server"`
	S3AccessKey  string        `required:"false" arg:"s3-access-key" env:"S3_ACCESS_KEY" usage:"Access Key for S3 server"`
	S3SecretKey  string        `required:"false" arg:"s3-secret-key" env:"S3_SECRET_KEY" usage:"Secret Key for S3 server" display:"length"`
	BucketName   string        `required:"true" arg:"bucket" env:"BUCKET" usage:"bucket to store backups in"`
}

func (a *application) Run(
	ctx context.Context,
	sentryClient libsentry.Client,
) error {
	s3Client := s3utils.CreateS3Client(s3utils.URL(a.S3Url), s3utils.AccessKey(a.S3AccessKey), s3utils.SecretKey(a.S3SecretKey))

	{
		glog.V(2).Infof("list all objects started")
		input := &s3.ListObjectsV2Input{
			Bucket: aws.String(a.BucketName),
		}
		paginator := s3.NewListObjectsV2Paginator(s3Client, input)

		for paginator.HasMorePages() {
			output, err := paginator.NextPage(ctx)
			if err != nil {
				return errors.Wrap(ctx, err, "list objects failed")
			}
			for _, object := range output.Contents {
				fmt.Printf("key: '%s'\n", *object.Key)
			}
		}
		glog.V(2).Infof("list all objects completed")
	}

	{
		glog.V(2).Infof("list prefix objects started")
		input := &s3.ListObjectsV2Input{
			Bucket: aws.String(a.BucketName),
			Prefix: aws.String("test"),
		}
		paginator := s3.NewListObjectsV2Paginator(s3Client, input)

		for paginator.HasMorePages() {
			output, err := paginator.NextPage(ctx)
			if err != nil {
				return errors.Wrap(ctx, err, "list objects failed")
			}
			for _, object := range output.Contents {
				fmt.Printf("key: '%s'\n", *object.Key)
			}
		}
		glog.V(2).Infof("list prefix objects completed")
	}

	return nil
}
