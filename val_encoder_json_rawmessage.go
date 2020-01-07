package jzon

import (
	"unsafe"
)

type jsonRawMessageEncoder struct{}

func (*jsonRawMessageEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	panic("not implemented")
}

func (*jsonRawMessageEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	data := *(*[]byte)(ptr)
	// TODO: raw message validation?
	s.Raw(data)
}
