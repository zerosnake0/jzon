package jzon

import (
	"unsafe"
)

// int8 encoder
type int8Encoder struct{}

func (*int8Encoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if ptr == nil {
		s.Null()
		return
	}
	s.Int8(*(*int8)(ptr))
}

// int16 encoder
type int16Encoder struct{}

func (*int16Encoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if ptr == nil {
		s.Null()
		return
	}
	s.Int16(*(*int16)(ptr))
}

// int32 encoder
type int32Encoder struct{}

func (*int32Encoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if ptr == nil {
		s.Null()
		return
	}
	s.Int32(*(*int32)(ptr))
}

// int64 encoder
type int64Encoder struct{}

func (*int64Encoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if ptr == nil {
		s.Null()
		return
	}
	s.Int64(*(*int64)(ptr))
}
