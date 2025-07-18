// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bborbe/errors"
	libsentry "github.com/bborbe/sentry"
	"github.com/bborbe/service"
	"github.com/golang/glog"

	libs3 "github.com/bborbe/s3-utils"
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
	S3Url        string        `required:"true" arg:"s3-url" env:"S3_URL" usage:"URL of S3 server"`
	S3AccessKey  string        `required:"true" arg:"s3-access-key" env:"S3_ACCESS_KEY" usage:"Access Key for S3 server"`
	S3SecretKey  string        `required:"true" arg:"s3-secret-key" env:"S3_SECRET_KEY" usage:"Secret Key for S3 server" display:"length"`
	BucketName   string        `required:"true" arg:"bucket" env:"BUCKET" usage:"bucket to store backups in"`
	KeyName      string        `required:"true" arg:"key" env:"KEY" usage:"key to store backups in"`
}

func (a *application) Run(
	ctx context.Context,
	sentryClient libsentry.Client,
) error {
	s3Client := libs3.CreateS3Client(a.S3Url, a.S3AccessKey, a.S3SecretKey)
	getObject, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(a.BucketName),
		Key:    aws.String(a.KeyName),
	})
	if err != nil {
		return errors.Wrap(ctx, err, "get object failed")
	}
	defer getObject.Body.Close()
	if _, err = io.Copy(os.Stdout, getObject.Body); err != nil {
		return errors.Wrapf(ctx, err, "download content failed")
	}
	glog.V(2).Infof("download %s completed", a.KeyName)
	return nil
}
