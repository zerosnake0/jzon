package jzon

import (
	"unsafe"
)

type jsonRawMessageEncoder struct{}

func (*jsonRawMessageEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return len(*(*[]byte)(ptr)) == 0
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
