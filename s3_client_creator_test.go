// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package s3utils_test

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	s3utils "github.com/bborbe/s3-utils"
)

var _ = Describe("CreateS3Client", func() {
	It("returns a non-nil client for http url", func() {
		client := s3utils.CreateS3Client("http://localhost:9000", "access", "secret")
		Expect(client).NotTo(BeNil())
	})

	It("returns a non-nil client for https url", func() {
		client := s3utils.CreateS3Client("https://s3.example.com", "access", "secret")
		Expect(client).NotTo(BeNil())
	})

	It("returns a non-nil client with empty credentials", func() {
		client := s3utils.CreateS3Client("http://localhost:9000", "", "")
		Expect(client).NotTo(BeNil())
	})

	It("signs requests with region garage in the credential scope", func() {
		client := s3utils.CreateS3Client("https://s3.example.com", "access", "secret")
		presigner := s3.NewPresignClient(client)
		req, err := presigner.PresignGetObject(context.Background(), &s3.GetObjectInput{
			Bucket: aws.String("b"), Key: aws.String("k"),
		}, func(o *s3.PresignOptions) {})
		Expect(err).NotTo(HaveOccurred())
		Expect(req.URL).To(ContainSubstring("%2Fgarage%2Fs3%2Faws4_request"))
	})
})
