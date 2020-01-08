package jzon

import (
	"unsafe"
)

// int8 encoder
type int8Encoder struct{}

func (*int8Encoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *(*int8)(ptr) == 0
}

func (*int8Encoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	if opts == nil || !opts.Quoted {
		s.Int8(*(*int8)(ptr))
	} else {
		s.quotedInt8(*(*int8)(ptr))
	}
}

// int16 encoder
type int16Encoder struct{}

func (*int16Encoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *(*int16)(ptr) == 0
}

func (*int16Encoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	if opts == nil || !opts.Quoted {
		s.Int16(*(*int16)(ptr))
	} else {
		s.quotedInt16(*(*int16)(ptr))
	}
}

// int32 encoder
type int32Encoder struct{}

func (*int32Encoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *(*int32)(ptr) == 0
}

func (*int32Encoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	if opts == nil || !opts.Quoted {
		s.Int32(*(*int32)(ptr))
	} else {
		s.quotedInt32(*(*int32)(ptr))
	}
}

// int64 encoder
type int64Encoder struct{}

func (*int64Encoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *(*int64)(ptr) == 0
}

func (*int64Encoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	if opts == nil || !opts.Quoted {
		s.Int64(*(*int64)(ptr))
	} else {
		s.quotedInt64(*(*int64)(ptr))
	}
}
