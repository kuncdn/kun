/*
Copyright 2019 The Koala Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package responsewriter

import (
	"bufio"
	"bytes"
	"net"
	"net/http"
)

// ResponseWriter .
type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	http.Hijacker
	http.CloseNotifier

	Reset()
	Body() []byte
	BodyString() string
	Written() bool
	Success() bool
}

// responseWriter is a ResponseWriter wrapper that may be used as buffer.
//
// A middleware may pass it to the next handlers ServeHTTP method as a
// drop in replacement for the response writer. After the ServeHTTP method is run the middleware may
// examine what has been written to the responseWriter and decide what to write to the "original" ResponseWriter
// (that may well be another buffer passed from another middleware).
//
// The downside is the body being written two times and the complete caching of the
// body in the memory which will be inacceptable for large bodies.
// Therefor Peek is an alternative response writer wrapper that only caching headers and status code
// but allowing to intercept calls of the Write method.
type responseWriter struct {

	// ResponseWriter is the underlying response writer that is wrapped by responseWriter
	http.ResponseWriter
	http.Flusher

	// if the underlying ResponseWriter is a Contexter, that Contexter is saved here

	// responseWriter is the underlying io.Writer that buffers the response body
	responseWriter bytes.Buffer

	// Code is the cached status code
	status int

	// changed tracks if anything has been set on the responsewriter. Also reads from the header
	// are seen as changes
	changed bool

	// header is the cached header
	header http.Header
}

var _ ResponseWriter = &responseWriter{}

// make sure to fulfill the Contexter interface

// NewResponseWriter creates a new responseWriter by wrapping the given response writer.
func NewResponseWriter(w http.ResponseWriter) ResponseWriter {
	bf := &responseWriter{}
	bf.ResponseWriter = w
	bf.header = make(http.Header)
	return bf
}

// Header returns the cached http.Header and tracks this call as change
func (bf *responseWriter) Header() http.Header {
	bf.changed = true
	return bf.header
}

// WriteHeader writes the cached status code and tracks this call as change
func (bf *responseWriter) WriteHeader(i int) {
	bf.changed = true
	bf.status = i
}

// Write writes to the underlying buffer and tracks this call as change
func (bf *responseWriter) Write(b []byte) (int, error) {
	bf.changed = true
	return bf.responseWriter.Write(b)
}

// Reset set the responseWriter to the defaults
func (bf *responseWriter) Reset() {
	bf.responseWriter.Reset()
	bf.status = 0
	bf.changed = false
	bf.header = make(http.Header)
}

// Flush flushes headers, status code and body to the underlying ResponseWriter, if something changed
func (bf *responseWriter) Flush() {
	if bf.Written() {
		if bf.status != 0 {
			bf.ResponseWriter.WriteHeader(bf.status)
		}
		header := bf.ResponseWriter.Header()
		for k, v := range bf.header {
			header.Del(k)
			for _, val := range v {
				header.Add(k, val)
			}
		}
		bf.ResponseWriter.Write(bf.responseWriter.Bytes())
	}
}

// Body returns the bytes of the underlying buffer (that is meant to be the body of the response)
func (bf *responseWriter) Body() []byte {
	return bf.responseWriter.Bytes()
}

// BodyString returns the string of the underlying buffer (that is meant to be the body of the response)
func (bf *responseWriter) BodyString() string {
	return bf.responseWriter.String()
}

// HasChanged returns true if Header, WriteHeader or Write has been called
func (bf *responseWriter) Written() bool {
	return bf.changed
}

// IsOk returns true if the cached status code is not set or in the 2xx range.
func (bf *responseWriter) Success() bool {
	if bf.status < http.StatusOK || bf.status > http.StatusBadRequest {
		return false
	}
	return true
}

func (bf *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return bf.ResponseWriter.(http.Hijacker).Hijack()
}

func (bf *responseWriter) CloseNotify() <-chan bool {
	return bf.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
