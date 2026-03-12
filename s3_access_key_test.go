// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package s3utils_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	s3utils "github.com/bborbe/s3-utils"
)

var _ = Describe("AccessKey", func() {
	It("String returns string value", func() {
		Expect(s3utils.AccessKey("myaccesskey").String()).To(Equal("myaccesskey"))
	})
	It("empty string", func() {
		Expect(s3utils.AccessKey("").String()).To(Equal(""))
	})
})
