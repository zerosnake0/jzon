package jzon

import (
	"unsafe"
)

// uint8 encoder
type uint8Encoder struct{}

func (*uint8Encoder) IsEmpty(ptr unsafe.Pointer) bool {
	if ptr == nil {
		return true
	}
	return *(*uint8)(ptr) == 0
}

func (*uint8Encoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	if opts == nil || !opts.Quoted {
		s.Uint8(*(*uint8)(ptr))
	} else {
		s.quotedUint8(*(*uint8)(ptr))
	}
}

// uint16 encoder
type uint16Encoder struct{}

func (*uint16Encoder) IsEmpty(ptr unsafe.Pointer) bool {
	if ptr == nil {
		return true
	}
	return *(*uint16)(ptr) == 0
}

func (*uint16Encoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	if opts == nil || !opts.Quoted {
		s.Uint16(*(*uint16)(ptr))
	} else {
		s.quotedUint16(*(*uint16)(ptr))
	}
}

// uint32 encoder
type uint32Encoder struct{}

func (*uint32Encoder) IsEmpty(ptr unsafe.Pointer) bool {
	if ptr == nil {
		return true
	}
	return *(*uint32)(ptr) == 0
}

func (*uint32Encoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	if opts == nil || !opts.Quoted {
		s.Uint32(*(*uint32)(ptr))
	} else {
		s.quotedUint32(*(*uint32)(ptr))
	}
}

// uint64 encoder
type uint64Encoder struct{}

func (*uint64Encoder) IsEmpty(ptr unsafe.Pointer) bool {
	if ptr == nil {
		return true
	}
	return *(*uint64)(ptr) == 0
}

func (*uint64Encoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	if opts == nil || !opts.Quoted {
		s.Uint64(*(*uint64)(ptr))
	} else {
		s.quotedUint64(*(*uint64)(ptr))
	}
}
