// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"os"
	"time"

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
}

func (a *application) Run(
	ctx context.Context,
	sentryClient libsentry.Client,
) error {
	s3Client := s3utils.CreateS3Client(s3utils.URL(a.S3Url), s3utils.AccessKey(a.S3AccessKey), s3utils.SecretKey(a.S3SecretKey))

	glog.V(2).Infof("list all objects started")
	buckets, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return errors.Wrap(ctx, err, "list buckets failed")
	}
	for _, bucket := range buckets.Buckets {
		fmt.Printf("name: '%s'\n", *bucket.Name)
	}
	glog.V(2).Infof("list all objects completed")

	return nil
}
