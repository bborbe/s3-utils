// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("copy", func() {
	It("compiles", func() {
		_, err := gexec.Build("github.com/bborbe/s3-utils/cmd/copy", "-mod=mod")
		Expect(err).NotTo(HaveOccurred())
	})

	It("copies all objects from source to destination", func() {
		src := newS3Stub("src")
		defer src.close()
		dst := newS3Stub("dst")
		defer dst.close()
		src.set("a", []byte("hello"))
		src.set("b", []byte("world"))
		src.set("c", []byte("garage"))

		Expect(runCopy(src, dst, false)).To(Succeed())
		for _, key := range []string{"a", "b", "c"} {
			got, ok := dst.get(key)
			want, _ := src.get(key)
			Expect(ok).To(BeTrue(), "key %q missing in destination", key)
			Expect(got).To(Equal(want), "key %q content mismatch", key)
		}
	})

	It("skips existing objects with matching size", func() {
		src := newS3Stub("src")
		defer src.close()
		dst := newS3Stub("dst")
		defer dst.close()
		src.set("same", []byte("12345")) // 5 bytes
		src.set("new", []byte("new-data"))
		dst.set("same", []byte("54321")) // same size (5 bytes), different content

		Expect(runCopy(src, dst, true)).To(Succeed())
		got, ok := dst.get("same")
		Expect(ok).To(BeTrue())
		Expect(
			string(got),
		).To(Equal("54321"), "matching-size object should be skipped, not overwritten")
		Expect(dst.putCountFor("same")).To(Equal(0))
		_, ok = dst.get("new")
		Expect(ok).To(BeTrue(), "new object missing in destination")
	})

	It("overwrites existing objects with different size", func() {
		src := newS3Stub("src")
		defer src.close()
		dst := newS3Stub("dst")
		defer dst.close()
		src.set("a", []byte("12345")) // 5 bytes
		dst.set("a", []byte("12"))    // 2 bytes — different size, must be overwritten

		Expect(runCopy(src, dst, true)).To(Succeed())
		got, _ := dst.get("a")
		Expect(string(got)).To(Equal("12345"), "different-size object should be overwritten")
	})

	It("retries transient GetObject failures", func() {
		src := newS3Stub("src")
		defer src.close()
		dst := newS3Stub("dst")
		defer dst.close()
		src.set("a", []byte("retry-me"))
		src.setFailGet("a", 1) // fail once, then succeed

		Expect(runCopy(src, dst, false)).To(Succeed())
		got, ok := dst.get("a")
		Expect(ok).To(BeTrue())
		Expect(string(got)).To(Equal("retry-me"))
	})

	It("surfaces HeadObject errors instead of treating them as absent", func() {
		src := newS3Stub("src")
		defer src.close()
		dst := newS3Stub("dst")
		defer dst.close()
		src.set("a", []byte("data"))
		dst.setHeadCode(http.StatusForbidden) // 403 must be surfaced

		Expect(runCopy(src, dst, true)).To(HaveOccurred())
		Expect(dst.putCountFor("a")).To(Equal(0), "no upload expected on HeadObject 403")
	})

	It("returns an error when listing the source fails", func() {
		src := newS3Stub("src")
		defer src.close()
		dst := newS3Stub("dst")
		defer dst.close()
		src.server.Close() // closing the server makes all requests fail

		Expect(runCopy(src, dst, false)).To(HaveOccurred())
	})
})

func TestCopySuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Copy Suite")
}
