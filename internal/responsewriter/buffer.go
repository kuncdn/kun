/*
Copyright 2019 The Tracfox Authors.
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
	"bytes"
	"net/http"
)

// Buffer is a ResponseWriter wrapper that may be used as buffer.
//
// A middleware may pass it to the next handlers ServeHTTP method as a
// drop in replacement for the response writer. After the ServeHTTP method is run the middleware may
// examine what has been written to the Buffer and decide what to write to the "original" ResponseWriter
// (that may well be another buffer passed from another middleware).
//
// The downside is the body being written two times and the complete caching of the
// body in the memory which will be inacceptable for large bodies.
// Therefor Peek is an alternative response writer wrapper that only caching headers and status code
// but allowing to intercept calls of the Write method.
type Buffer struct {

	// ResponseWriter is the underlying response writer that is wrapped by Buffer
	http.ResponseWriter

	// if the underlying ResponseWriter is a Contexter, that Contexter is saved here

	// Buffer is the underlying io.Writer that buffers the response body
	Buffer bytes.Buffer

	// Code is the cached status code
	Code int

	// changed tracks if anything has been set on the responsewriter. Also reads from the header
	// are seen as changes
	changed bool

	// header is the cached header
	header http.Header
}

// make sure to fulfill the Contexter interface

// NewBuffer creates a new Buffer by wrapping the given response writer.
func NewBuffer(w http.ResponseWriter) (bf *Buffer) {
	bf = &Buffer{}
	bf.ResponseWriter = w
	bf.header = make(http.Header)
	return
}

// Header returns the cached http.Header and tracks this call as change
func (bf *Buffer) Header() http.Header {
	bf.changed = true
	return bf.header
}

// WriteHeader writes the cached status code and tracks this call as change
func (bf *Buffer) WriteHeader(i int) {
	bf.changed = true
	bf.Code = i
}

// Write writes to the underlying buffer and tracks this call as change
func (bf *Buffer) Write(b []byte) (int, error) {
	bf.changed = true
	return bf.Buffer.Write(b)
}

// Reset set the Buffer to the defaults
func (bf *Buffer) Reset() {
	bf.Buffer.Reset()
	bf.Code = 0
	bf.changed = false
	bf.header = make(http.Header)
}

// FlushAll flushes headers, status code and body to the underlying ResponseWriter, if something changed
func (bf *Buffer) FlushAll() {
	if bf.HasChanged() {
		bf.FlushHeaders()
		bf.FlushCode()
		bf.ResponseWriter.Write(bf.Buffer.Bytes())
	}
}

// Body returns the bytes of the underlying buffer (that is meant to be the body of the response)
func (bf *Buffer) Body() []byte {
	return bf.Buffer.Bytes()
}

// BodyString returns the string of the underlying buffer (that is meant to be the body of the response)
func (bf *Buffer) BodyString() string {
	return bf.Buffer.String()
}

// HasChanged returns true if Header, WriteHeader or Write has been called
func (bf *Buffer) HasChanged() bool {
	return bf.changed
}

// IsOk returns true if the cached status code is not set or in the 2xx range.
func (bf *Buffer) IsOk() bool {
	if bf.Code == 0 {
		return true
	}
	if bf.Code >= 200 && bf.Code < 300 {
		return true
	}
	return false
}

// FlushCode flushes the status code to the underlying responsewriter if it was set.
func (bf *Buffer) FlushCode() {
	if bf.Code != 0 {
		bf.ResponseWriter.WriteHeader(bf.Code)
	}
}

// FlushHeaders adds the headers to the underlying ResponseWriter, removing them from Buffer.
func (bf *Buffer) FlushHeaders() {
	header := bf.ResponseWriter.Header()
	for k, v := range bf.header {
		header.Del(k)
		for _, val := range v {
			header.Add(k, val)
		}
	}
}
