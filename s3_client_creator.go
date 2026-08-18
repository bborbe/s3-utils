// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package s3utils

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
	"github.com/golang/glog"
)

// CreateS3ClientOption defines a function for modifying the s3.Options
// of a client created by CreateS3Client.
type CreateS3ClientOption func(*s3.Options)

func CreateS3Client(
	s3Url URL,
	s3AccessKey AccessKey,
	s3SecretKey SecretKey,
	opts ...CreateS3ClientOption,
) *s3.Client {
	options := s3.Options{
		UsePathStyle:     true,
		EndpointResolver: s3.EndpointResolverFromURL(s3Url.String()),
		EndpointOptions: s3.EndpointResolverOptions{
			DisableHTTPS: !strings.HasPrefix(s3Url.String(), "https://"),
		},
		Credentials: aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider(
				s3AccessKey.String(),
				s3SecretKey.String(),
				"",
			),
		),
	}
	for _, opt := range opts {
		opt(&options)
	}
	return s3.New(
		options,
		func(o *s3.Options) {
			if glog.V(4) {
				o.Logger = logging.StandardLogger{
					Logger: log.Default(),
				}
				o.ClientLogMode = aws.LogSigning |
					aws.LogRetries |
					aws.LogRequest |
					aws.LogRequestWithBody |
					aws.LogResponse |
					aws.LogResponseWithBody |
					aws.LogDeprecatedUsage |
					aws.LogRequestEventMessage |
					aws.LogResponseEventMessage

			}
		},
	)
}

// WithRegion sets the region used for SigV4 signing. An empty region is the
// SDK default; callers targeting Garage should pass WithRegion("garage").
func WithRegion(region string) CreateS3ClientOption {
	return func(options *s3.Options) {
		options.Region = region
	}
}
