package jzon

import (
	"unsafe"
)

type efaceEncoder struct{}

func (*efaceEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *(*interface{})(ptr) == nil
}

func (*efaceEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	s.value(*(*interface{})(ptr), opts)
}

type ifaceEncoder struct{}

func (*ifaceEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	// TODO: is this ok?
	o := packIFace(ptr)
	return o == nil
}

func (*ifaceEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	o := packIFace(ptr)
	s.value(o, opts)
}
