package jzon

import (
	"unsafe"
)

// bool encoder
type boolEncoder struct{}

func (*boolEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if ptr == nil {
		s.Null()
		return
	}
	s.Bool(*(*bool)(ptr))
}

// string encoder
type stringEncoder struct{}

func (*stringEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if ptr == nil {
		s.Null()
		return
	}
	s.String(*(*string)(ptr))
}
