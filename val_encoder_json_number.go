package jzon

import (
	"unsafe"
)

type jsonNumberEncoder struct {
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
	// TODO: the standard lib will check the validity
	s.RawString(str)
}
