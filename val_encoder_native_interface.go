package jzon

import (
	"unsafe"
)

type efaceEncoder struct {
}

func (enc *efaceEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	s.Value(*(*interface{})(ptr))
}

type ifaceEncoder struct {
}

func (enc *ifaceEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	o := packIFace(ptr)
	s.Value(o)
}
