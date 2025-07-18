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

type URL string

func (u URL) String() string {
	return string(u)
}

type AccessKey string

func (a AccessKey) String() string {
	return string(a)
}

type SecretKey string

func (s SecretKey) String() string {
	return string(s)
}

func CreateS3Client(
	s3Url URL,
	s3AccessKey AccessKey,
	s3SecretKey SecretKey,
) *s3.Client {
	return s3.New(
		s3.Options{
			UsePathStyle:     true,
			EndpointResolver: s3.EndpointResolverFromURL(s3Url.String()),
			EndpointOptions: s3.EndpointResolverOptions{
				DisableHTTPS: !strings.HasPrefix("https://", s3Url.String()),
			},
			Credentials: aws.NewCredentialsCache(
				credentials.NewStaticCredentialsProvider(
					s3AccessKey.String(),
					s3SecretKey.String(),
					"",
				),
			),
		},
		func(options *s3.Options) {
			if glog.V(4) {
				options.Logger = logging.StandardLogger{
					Logger: log.Default(),
				}
				options.ClientLogMode = aws.LogSigning |
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
