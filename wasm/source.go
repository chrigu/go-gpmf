package main

import (
	"errors"
	"io"
	"syscall/js"
)

const page = 64 * 1024 // size of our scratch buffer

type Source struct {
	size    int64    // total length of the resource
	pos     int64    // current cursor
	getBuf  js.Value // JS function (offset, length) → Uint8Array
	scratch []byte   // reusable copy target
}

// NewSource builds a ReadSeeker. getBuf **must** be synchronous
// (e.g. use a pre-fetched ArrayBuffer or an in-memory cache).
func NewSource(size int64, getBuf js.Value) *Source {
	return &Source{
		size:    size,
		getBuf:  getBuf,
		scratch: make([]byte, page),
	}
}

// --- io.Reader ---
func (s *Source) Read(p []byte) (int, error) {
	if s.pos >= s.size {
		return 0, io.EOF
	}
	want := int64(len(p))
	if remain := s.size - s.pos; want > remain {
		want = remain
	}
	// JS: buf := getBuf(offset, length) ➜ Uint8Array
	view := s.getBuf.Invoke(s.pos, want)
	if !view.InstanceOf(js.Global().Get("Uint8Array")) {
		return 0, errors.New("getBuf did not return Uint8Array")
	}
	n := js.CopyBytesToGo(p, view)
	s.pos += int64(n)
	return n, nil
}

// --- io.Seeker ---
func (s *Source) Seek(off int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		s.pos = off
	case io.SeekCurrent:
		s.pos += off
	case io.SeekEnd:
		s.pos = s.size + off
	default:
		return 0, errors.New("invalid whence")
	}
	if s.pos < 0 {
		s.pos = 0
	}
	if s.pos > s.size {
		s.pos = s.size
	}
	return s.pos, nil
}
