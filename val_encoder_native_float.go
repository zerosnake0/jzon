package jzon

import (
	"unsafe"
)

// float32 encoder
type float32Encoder struct{}

func (*float32Encoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if ptr == nil {
		s.Null()
		return
	}
	s.Float32(*(*float32)(ptr))
}

// float64 encoder
type float64Encoder struct{}

func (*float64Encoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if ptr == nil {
		s.Null()
		return
	}
	s.Float64(*(*float64)(ptr))
}