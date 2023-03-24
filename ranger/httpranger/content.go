// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the Go project's LICENSE file.
//
// This file was lifted wholesale from the Go standard library with
// very minor tweaks by Storx Labs, Inc., 2018

package httpranger

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/zeebo/errs"

	"common/ranger"
)

// ErrInvalidRange is an error returned when an invalid HTTP range is requested.
var ErrInvalidRange = errs.Class("invalid range")

// ServeContent is the Go standard library's http.ServeContent but modified to
// work with Rangers.
func ServeContent(ctx context.Context, w http.ResponseWriter, r *http.Request, name string, modtime time.Time, content ranger.Ranger) (err error) {
	defer mon.Task()(&ctx)(&err)

	// If the Content-Type is specified, use it, otherwise do the cheap
	// heuristic using the suffix lookup.
	ctypes, haveType := w.Header()["Content-Type"]

	var ctype string
	if !haveType {
		ctype = mime.TypeByExtension(filepath.Ext(name))
	} else if len(ctypes) > 0 {
		ctype = ctypes[0]
	}

	if ctype != "" {
		w.Header().Set("Content-Type", ctype)
	}

	setLastModified(w, modtime)
	done, rangeReq := checkPreconditions(w, r, modtime)
	if done {
		return nil
	}

	code := http.StatusOK

	size := content.Size()

	if size <= 0 {
		w.WriteHeader(code)
		return nil
	}

	// handle Content-Range header.
	sendSize := size
	sendContent := func() (io.ReadCloser, error) {
		return content.Range(ctx, 0, size)
	}

	ranges, err := ParseRange(rangeReq, size)
	if err != nil {
		if errors.Is(err, errNoOverlap) {
			w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", size))
		}
		return err
	}
	if sumRangesSize(ranges) > size {
		// The total number of bytes in all the ranges
		// is larger than the size of the file by
		// itself, so this is probably an attack, or a
		// dumb client. Ignore the range request.
		ranges = nil
	}
	switch {
	case len(ranges) == 1:
		// RFC 2616, Section 14.16:
		// "When an HTTP message includes the content of a single
		// range (for example, a response to a request for a
		// single range, or to a request for a set of ranges
		// that overlap without any holes), this content is
		// transmitted with a Content-Range header, and a
		// Content-Length header showing the number of bytes
		// actually transferred.
		// ...
		// A response to a request for a single range MUST NOT
		// be sent using the multipart/byteranges media type."
		ra := ranges[0]
		sendContent = func() (io.ReadCloser, error) { return content.Range(ctx, ra.Start, ra.Length) }
		sendSize = ra.Length
		code = http.StatusPartialContent
		w.Header().Set("Content-Range", ra.contentRange(size))
	case len(ranges) > 1:
		sendSize = rangesMIMESize(ranges, ctype, size)
		code = http.StatusPartialContent

		pr, pw := io.Pipe()
		mw := multipart.NewWriter(pw)
		w.Header().Set("Content-Type",
			"multipart/byteranges; boundary="+mw.Boundary())
		sendContent = func() (io.ReadCloser, error) { return io.NopCloser(pr), nil }
		// cause writing goroutine to fail and exit if CopyN doesn't finish.
		defer func() {
			if err := pr.Close(); err != nil {
				log.Printf("Error Closing pipereader: %s", err)
			}
		}()

		go func() {
			for _, ra := range ranges {
				part, err := mw.CreatePart(ra.mimeHeader(ctype, size))
				if err != nil {
					if err := pw.CloseWithError(err); err != nil {
						log.Printf("Error Closing pipewriter with errors: %s", err)
					}

					return
				}
				partReader, err := content.Range(ctx, ra.Start, ra.Length)
				if err != nil {
					if err := pw.CloseWithError(err); err != nil {
						log.Printf("Error Closing pipewriter with errors: %s", err)
					}

					return
				}
				defer func() {
					if err := partReader.Close(); err != nil {
						log.Printf("Error Closing partReader: %s", err)
					}
				}()

				if _, err := io.Copy(part, partReader); err != nil {
					if err := pw.CloseWithError(err); err != nil {
						log.Printf("Error Closing pipewriter with errors: %s", err)
					}

					return
				}
			}
			if err := mw.Close(); err != nil {
				log.Printf("Error Closing writer: %s", err)
			}

			if err := pw.Close(); err != nil {
				log.Printf("Error closing pipewriter: %s", err)
			}
		}()
	}

	w.Header().Set("Accept-Ranges", "bytes")
	if w.Header().Get("Content-Encoding") == "" {
		w.Header().Set("Content-Length", strconv.FormatInt(sendSize, 10))
	}

	if r.Method == http.MethodHead {
		w.WriteHeader(code)
		return nil
	}

	rd, err := sendContent()
	if err != nil {
		delete(w.Header(), "Content-Type")
		delete(w.Header(), "Content-Length")
		delete(w.Header(), "Last-Modified")
		return err
	}

	w.WriteHeader(code)

	if _, err := io.CopyN(w, rd, sendSize); err != nil {
		log.Printf("Error Copying bytes: %s", err)
	}

	if err := rd.Close(); err != nil {
		log.Printf("Error closing: %s", err)
	}

	return nil
}

var unixEpochTime = time.Unix(0, 0)

// isZeroTime reports whether t is obviously unspecified (either zero or
// Unix()=0).
func isZeroTime(t time.Time) bool {
	return t.IsZero() || t.Equal(unixEpochTime)
}

func setLastModified(w http.ResponseWriter, modtime time.Time) {
	if w == nil {
		return
	}

	if !isZeroTime(modtime) {
		w.Header().Set("Last-Modified", modtime.UTC().Format(http.TimeFormat))
	}
}

// checkPreconditions evaluates request preconditions and reports whether a
// precondition resulted in sending StatusNotModified or
// StatusPreconditionFailed.
func checkPreconditions(w http.ResponseWriter, r *http.Request, modtime time.Time) (done bool, rangeHeader string) {
	// This function carefully follows RFC 7232 section 6.
	ch := checkIfMatch(w, r)
	if ch == condNone {
		ch = checkIfUnmodifiedSince(r, modtime)
	}
	if ch == condFalse {
		w.WriteHeader(http.StatusPreconditionFailed)
		return true, ""
	}
	switch checkIfNoneMatch(w, r) {
	case condFalse:
		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			writeNotModified(w)
			return true, ""
		}
		w.WriteHeader(http.StatusPreconditionFailed)
		return true, ""
	case condNone:
		if checkIfModifiedSince(r, modtime) == condFalse {
			writeNotModified(w)
			return true, ""
		}
	}

	rangeHeader = r.Header.Get("Range")
	if rangeHeader != "" && checkIfRange(w, r, modtime) == condFalse {
		rangeHeader = ""
	}
	return false, rangeHeader
}

// condResult is the result of an HTTP request precondition check.
// See https://tools.ietf.org/html/rfc7232 section 3.
type condResult int

const (
	condNone condResult = iota
	condTrue
	condFalse
)

func checkIfMatch(w http.ResponseWriter, r *http.Request) condResult {
	im := r.Header.Get("If-Match")
	if im == "" {
		return condNone
	}
	for {
		im = textproto.TrimString(im)
		if len(im) == 0 {
			break
		}
		if im[0] == ',' {
			im = im[1:]
			continue
		}
		if im[0] == '*' {
			return condTrue
		}
		etag, remain := scanETag(im)
		if etag == "" {
			break
		}
		if etagStrongMatch(etag, w.Header().Get("Etag")) {
			return condTrue
		}
		im = remain
	}

	return condFalse
}

func checkIfUnmodifiedSince(r *http.Request, modtime time.Time) condResult {
	ius := r.Header.Get("If-Unmodified-Since")
	if ius == "" || isZeroTime(modtime) {
		return condNone
	}
	if t, err := http.ParseTime(ius); err == nil {
		// The Date-Modified header truncates sub-second precision, so
		// use mtime < t+1s instead of mtime <= t to check for unmodified.
		if modtime.Before(t.Add(1 * time.Second)) {
			return condTrue
		}
		return condFalse
	}
	return condNone
}

func checkIfNoneMatch(w http.ResponseWriter, r *http.Request) condResult {
	inm := r.Header.Get("If-None-Match")
	if inm == "" {
		return condNone
	}
	buf := inm
	for {
		buf = textproto.TrimString(buf)
		if len(buf) == 0 {
			break
		}
		if buf[0] == ',' {
			buf = buf[1:]
			continue
		}
		if buf[0] == '*' {
			return condFalse
		}
		etag, remain := scanETag(buf)
		if etag == "" {
			break
		}
		if etagWeakMatch(etag, w.Header().Get("Etag")) {
			return condFalse
		}
		buf = remain
	}
	return condTrue
}

func checkIfModifiedSince(r *http.Request, modtime time.Time) condResult {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return condNone
	}
	ims := r.Header.Get("If-Modified-Since")
	if ims == "" || isZeroTime(modtime) {
		return condNone
	}
	t, err := http.ParseTime(ims)
	if err != nil {
		return condNone
	}
	// The Date-Modified header truncates sub-second precision, so
	// use mtime < t+1s instead of mtime <= t to check for unmodified.
	if modtime.Before(t.Add(1 * time.Second)) {
		return condFalse
	}
	return condTrue
}

func checkIfRange(w http.ResponseWriter, r *http.Request, modtime time.Time) (rv condResult) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return condNone
	}
	ir := r.Header.Get("If-Range")
	if ir == "" {
		return condNone
	}
	etag, _ := scanETag(ir)
	if etag != "" {
		if etagStrongMatch(etag, w.Header().Get("Etag")) {
			return condTrue
		}
		return condFalse
	}
	// The If-Range value is typically the ETag value, but it may also be
	// the modtime date. See golang.org/issue/8367.
	if isZeroTime(modtime) {
		return condFalse
	}
	t, err := http.ParseTime(ir)
	if err != nil {
		return condFalse
	}
	if t.Unix() == modtime.Unix() {
		return condTrue
	}
	return condFalse
}

func writeNotModified(w http.ResponseWriter) {
	// RFC 7232 section 4.1:
	// a sender SHOULD NOT generate representation metadata other than the
	// above listed fields unless said metadata exists for the purpose of
	// guiding cache updates (e.g., Last-Modified might be useful if the
	// response does not have an ETag field).
	h := w.Header()
	delete(h, "Content-Type")
	delete(h, "Content-Length")
	if h.Get("Etag") != "" {
		delete(h, "Last-Modified")
	}
	w.WriteHeader(http.StatusNotModified)
}

// scanETag determines if a syntactically valid ETag is present at s. If so,
// the ETag and remaining text after consuming ETag is returned. Otherwise,
// it returns "", "".
func scanETag(s string) (etag string, remain string) {
	s = textproto.TrimString(s)
	start := 0
	if strings.HasPrefix(s, "W/") {
		start = 2
	}
	if len(s[start:]) < 2 || s[start] != '"' {
		return "", ""
	}
	// ETag is either W/"text" or "text".
	// See RFC 7232 2.3.
	for i := start + 1; i < len(s); i++ {
		c := s[i]
		switch {
		// Character values allowed in ETags.
		case c == 0x21 || c >= 0x23 && c <= 0x7A || c >= 0x80:
		case c == '"':
			return s[:i+1], s[i+1:]
		default:
			return "", ""
		}
	}
	return "", ""
}

// etagStrongMatch reports whether a and b match using strong ETag comparison.
// Assumes a and b are valid ETags.
func etagStrongMatch(a, b string) bool {
	return a == b && a != "" && a[0] == '"'
}

// etagWeakMatch reports whether a and b match using weak ETag comparison.
// Assumes a and b are valid ETags.
func etagWeakMatch(a, b string) bool {
	return strings.TrimPrefix(a, "W/") == strings.TrimPrefix(b, "W/")
}

// HTTPRange specifies the byte range to be sent to the client.
type HTTPRange struct {
	Start, Length int64
}

func (r HTTPRange) contentRange(size int64) string {
	return fmt.Sprintf("bytes %d-%d/%d", r.Start, r.Start+r.Length-1, size)
}

func (r HTTPRange) mimeHeader(contentType string, size int64) (rv textproto.MIMEHeader) {
	return textproto.MIMEHeader{
		"Content-Range": {r.contentRange(size)},
		"Content-Type":  {contentType},
	}
}

// ParseRange parses a Range header string as per RFC 2616.
// errNoOverlap is returned if none of the ranges overlap.
func ParseRange(s string, size int64) ([]HTTPRange, error) {
	if s == "" {
		return nil, nil // header not present
	}
	const b = "bytes="
	if !strings.HasPrefix(s, b) {
		return nil, ErrInvalidRange.New(s)
	}

	var ranges []HTTPRange
	noOverlap := false
	for _, ra := range strings.Split(s[len(b):], ",") {
		ra = strings.TrimSpace(ra)
		if ra == "" {
			continue
		}
		i := strings.Index(ra, "-")
		if i < 0 {
			return nil, ErrInvalidRange.New(ra)
		}
		start, end := strings.TrimSpace(ra[:i]), strings.TrimSpace(ra[i+1:])
		var r HTTPRange
		if start == "" {
			// If no start is specified, end specifies the
			// range start relative to the end of the file,
			// and we are dealing with <suffix-length>
			// which has to be a non-negative integer as per
			// RFC 7233 Section 2.1 "Byte-Ranges".
			i, err := strconv.ParseInt(end, 10, 64)
			if err != nil {
				return nil, ErrInvalidRange.New(ra)
			}
			if i > size {
				i = size
			}
			r.Start = size - i
			r.Length = size - r.Start
		} else {
			i, err := strconv.ParseInt(start, 10, 64)
			if err != nil || i < 0 {
				return nil, ErrInvalidRange.New(ra)
			}
			if i >= size {
				// If the range begins after the size of the content,
				// then it does not overlap.
				noOverlap = true
				continue
			}
			r.Start = i
			if end == "" {
				// If no end is specified, range extends to end of the file.
				r.Length = size - r.Start
			} else {
				i, err := strconv.ParseInt(end, 10, 64)
				if err != nil || r.Start > i {
					return nil, ErrInvalidRange.New(ra)
				}
				if i >= size {
					i = size - 1
				}
				r.Length = i - r.Start + 1
			}
		}
		ranges = append(ranges, r)
	}
	if noOverlap && len(ranges) == 0 {
		// The specified ranges did not overlap with the content.
		return nil, errNoOverlap
	}
	return ranges, nil
}

// countingWriter counts how many bytes have been written to it.
type countingWriter int64

func (w *countingWriter) Write(p []byte) (n int, err error) {
	*w += countingWriter(len(p))
	return len(p), nil
}

// rangesMIMESize returns the number of bytes it takes to encode the
// provided ranges as a multipart response.
func rangesMIMESize(ranges []HTTPRange, contentType string, contentSize int64) (encSize int64) {
	var w countingWriter
	mw := multipart.NewWriter(&w)
	for _, ra := range ranges {
		if _, err := mw.CreatePart(ra.mimeHeader(contentType, contentSize)); err != nil {
			log.Printf("Failed to Create Part: %s", err)
		}

		encSize += ra.Length
	}
	if err := mw.Close(); err != nil {
		log.Printf("Failed to close: %s", err)
	}

	encSize += int64(w)
	return encSize
}

func sumRangesSize(ranges []HTTPRange) (size int64) {
	for _, ra := range ranges {
		size += ra.Length
	}
	return size
}

// errNoOverlap is returned by serveContent's ParseRange if first-byte-pos of
// all of the byte-range-spec values is greater than the content size.
var errNoOverlap = ErrInvalidRange.New("failed to overlap")
