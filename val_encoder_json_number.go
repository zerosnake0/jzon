package jzon

import (
	"unsafe"
)

type jsonNumberEncoder struct {
}

func (*jsonNumberEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *(*string)(ptr) == ""
}

func (*jsonNumberEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	str := *(*string)(ptr)
	if str == "" {
		str = "0"
	}
	// TODO: the standard library will check the validity in future
	s.RawString(str)
}
