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
	"time"

	"github.com/getsentry/sentry-go"
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
