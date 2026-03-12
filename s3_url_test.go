// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package s3utils_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	s3utils "github.com/bborbe/s3-utils"
)

var _ = Describe("URL", func() {
	It("String returns string value", func() {
		Expect(s3utils.URL("http://localhost:9000").String()).To(Equal("http://localhost:9000"))
	})
	It("empty string", func() {
		Expect(s3utils.URL("").String()).To(Equal(""))
	})
})
