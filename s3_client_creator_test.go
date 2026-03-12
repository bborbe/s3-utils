// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package s3utils_test

import (
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
})
