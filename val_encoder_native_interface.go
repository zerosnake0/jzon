package jzon

import (
	"unsafe"
)

type efaceEncoder struct{}

func (*efaceEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	s.value(*(*interface{})(ptr), opts)
}

type ifaceEncoder struct{}

func (*ifaceEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	o := packIFace(ptr)
	s.value(o, opts)
}
