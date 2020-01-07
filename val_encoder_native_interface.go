package jzon

import (
	"unsafe"
)

type efaceEncoder struct{}

func (*efaceEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	if ptr == nil {
		return true
	}
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
	panic("not implemented")
}

func (*ifaceEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	o := packIFace(ptr)
	s.value(o, opts)
}
