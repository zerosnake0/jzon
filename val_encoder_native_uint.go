package jzon

import (
	"unsafe"
)

// uint8 encoder
type uint8Encoder struct{}

func (*uint8Encoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if ptr == nil {
		s.Null()
		return
	}
	s.Uint8(*(*uint8)(ptr))
}

// uint16 encoder
type uint16Encoder struct{}

func (*uint16Encoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if ptr == nil {
		s.Null()
		return
	}
	s.Uint16(*(*uint16)(ptr))
}

// uint32 encoder
type uint32Encoder struct{}

func (*uint32Encoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if ptr == nil {
		s.Null()
		return
	}
	s.Uint32(*(*uint32)(ptr))
}

// uint64 encoder
type uint64Encoder struct{}

func (*uint64Encoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if ptr == nil {
		s.Null()
		return
	}
	s.Uint64(*(*uint64)(ptr))
}
