package jzon

import (
	"unsafe"
)

type boolEncoder struct{}

func (*boolEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if ptr == nil {
		s.Null()
		return
	}
	s.Bool(*(*bool)(ptr))
}
