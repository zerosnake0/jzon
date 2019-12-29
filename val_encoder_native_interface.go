package jzon

import (
	"unsafe"
)

type efaceEncoder struct {
}

func (enc *efaceEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	s.Value(*(*interface{})(ptr))
}

type ifaceEncoder struct {
}

func (enc *ifaceEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	o := packIFace(ptr)
	s.Value(o)
}
