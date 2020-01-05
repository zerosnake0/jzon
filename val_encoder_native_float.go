package jzon

import (
	"reflect"
	"unsafe"
)

// float32 encoder
type float32Encoder struct{}

func (*float32Encoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	s.Float32(*(*float32)(ptr))
}

func (*float32Encoder) Encode2(v reflect.Value, s *Streamer, opts *EncOpts) {
	ptr := ptrOfValue(v)
	s.Float32(*(*float32)(ptr))
}

// float64 encoder
type float64Encoder struct{}

func (*float64Encoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	s.Float64(*(*float64)(ptr))
}

func (*float64Encoder) Encode2(v reflect.Value, s *Streamer, opts *EncOpts) {
	ptr := ptrOfValue(v)
	s.Float64(*(*float64)(ptr))
}
