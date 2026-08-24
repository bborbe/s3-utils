// Copyright (c) 2025 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

// sentryStub satisfies the bborbe/sentry Client interface without network I/O.
type sentryStub struct{}

func (s sentryStub) CaptureMessage(
	string,
	*sentry.EventHint,
	sentry.EventModifier,
) *sentry.EventID {
	return nil
}

func (s sentryStub) CaptureException(
	error,
	*sentry.EventHint,
	sentry.EventModifier,
) *sentry.EventID {
	return nil
}

func (s sentryStub) Flush(time.Duration) bool { return true }

func (s sentryStub) Close() error { return nil }

// s3Stub is a minimal in-memory S3-compatible server implementing the
// operations the copy command exercises: ListObjectsV2, GetObject,
// PutObject, HeadObject. It intentionally fails GetObject for keys listed
// in failGet (simulating transient errors) and can force a HeadObject
// status via headStatus.
type s3Stub struct {
	mu       sync.Mutex
	objects  map[string][]byte
	failGet  map[string]int
	headCode int
	putCount map[string]int
	server   *httptest.Server
	bucket   string
}

func newS3Stub(bucket string) *s3Stub {
	s := &s3Stub{
		objects:  make(map[string][]byte),
		failGet:  make(map[string]int),
		putCount: make(map[string]int),
		bucket:   bucket,
		headCode: http.StatusOK,
	}
	s.server = httptest.NewServer(http.HandlerFunc(s.handle))
	return s
}

func (s *s3Stub) close() { s.server.Close() }

func (s *s3Stub) url() string { return s.server.URL }

// set sets an object in the bucket.
func (s *s3Stub) set(key string, data []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.objects[key] = data
}

// setFailGet makes the next n GetObject calls for key return 500.
func (s *s3Stub) setFailGet(key string, n int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.failGet[key] = n
}

// setHeadCode forces HeadObject to return the given status (e.g. 403).
func (s *s3Stub) setHeadCode(code int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.headCode = code
}

func (s *s3Stub) get(key string) ([]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, ok := s.objects[key]
	return data, ok
}

func (s *s3Stub) putCountFor(key string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.putCount[key]
}

func (s *s3Stub) handle(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	path := strings.TrimPrefix(r.URL.Path, "/")
	// Path-style addressing: /<bucket>/<key> or /<bucket> (list)
	rest := strings.TrimPrefix(path, s.bucket+"/")
	switch {
	case r.Method == http.MethodGet && r.URL.Query().Get("list-type") == "2":
		s.handleList(w)
	case r.Method == http.MethodGet:
		s.handleGet(w, rest)
	case r.Method == http.MethodPut:
		s.handlePut(w, r, rest)
	case r.Method == http.MethodHead:
		s.handleHead(w, rest)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (s *s3Stub) handleList(w http.ResponseWriter) {
	type content struct {
		Key          string
		LastModified string
		ETag         string
		Size         int64
		StorageClass string
	}
	type result struct {
		XMLName     xml.Name `xml:"ListBucketResult"`
		Ns          string   `xml:"xmlns,attr"`
		Name        string
		IsTruncated bool
		KeyCount    int
		MaxKeys     int
		Contents    []content
	}
	var contents []content
	for key, data := range s.objects {
		contents = append(contents, content{
			Key:          key,
			LastModified: "2026-08-24T00:00:00.000Z",
			ETag:         fmt.Sprintf("\"%d\"", len(data)),
			Size:         int64(len(data)),
			StorageClass: "STANDARD",
		})
	}
	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(result{
		Ns:          "http://s3.amazonaws.com/doc/2006-03-01/",
		Name:        s.bucket,
		IsTruncated: false,
		KeyCount:    len(contents),
		MaxKeys:     1000,
		Contents:    contents,
	})
}

func (s *s3Stub) handleGet(w http.ResponseWriter, key string) {
	if fails := s.failGet[key]; fails > 0 {
		s.failGet[key] = fails - 1
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, ok := s.objects[key]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write(
			[]byte(
				`<?xml version="1.0"?><Error><Code>NoSuchKey</Code><Message>not found</Message></Error>`,
			),
		)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
}

func (s *s3Stub) handlePut(w http.ResponseWriter, r *http.Request, key string) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.objects[key] = data
	s.putCount[key]++
	w.Header().Set("ETag", fmt.Sprintf("\"%d\"", len(data)))
	w.WriteHeader(http.StatusOK)
}

func (s *s3Stub) handleHead(w http.ResponseWriter, key string) {
	if s.headCode != http.StatusOK {
		w.WriteHeader(s.headCode)
		return
	}
	data, ok := s.objects[key]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
}

func runCopy(src, dst *s3Stub, skipExisting bool) error {
	app := &application{
		SrcS3Url:       src.url(),
		SrcS3AccessKey: "src-key",
		SrcS3SecretKey: "src-secret",
		SrcS3Region:    "garage",
		SrcBucket:      src.bucket,
		DstS3Url:       dst.url(),
		DstS3AccessKey: "dst-key",
		DstS3SecretKey: "dst-secret",
		DstS3Region:    "garage",
		DstBucket:      dst.bucket,
		Concurrency:    4,
		SkipExisting:   skipExisting,
	}
	return app.Run(context.Background(), sentryStub{})
}

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
