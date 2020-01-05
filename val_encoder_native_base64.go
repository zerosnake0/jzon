package jzon

import (
	"encoding/base64"
	"reflect"
	"unsafe"
)

type base64Encoder struct{}

func (*base64Encoder) Encode(ptr unsafe.Pointer, s *Streamer, _ *EncOpts) {
	if s.Error != nil {
		return
	}
	if ptr == nil {
		s.Null()
		return
	}
	src := *(*[]byte)(ptr)
	if src == nil {
		s.Null()
		return
	}
	s.onVal()
	s.buffer = append(s.buffer, '"')
	if l := len(src); l != 0 {
		enc := base64.StdEncoding
		size := enc.EncodedLen(l)
		buf := make([]byte, size)
		enc.Encode(buf, src)
		s.buffer = append(s.buffer, buf...)
	}
	s.buffer = append(s.buffer, '"')
}

func (*base64Encoder) Encode2(v reflect.Value, s *Streamer, _ *EncOpts) {
	if s.Error != nil {
		return
	}
	if v.IsNil() {
		s.Null()
		return
	}
	src := v.Bytes()
	if src == nil {
		s.Null()
		return
	}
	s.onVal()
	s.buffer = append(s.buffer, '"')
	if l := len(src); l != 0 {
		enc := base64.StdEncoding
		size := enc.EncodedLen(l)
		buf := make([]byte, size)
		enc.Encode(buf, src)
		s.buffer = append(s.buffer, buf...)
	}
	s.buffer = append(s.buffer, '"')
}
