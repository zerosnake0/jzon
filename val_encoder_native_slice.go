package jzon

import (
	"reflect"
	"unsafe"
)

type sliceEncoderBuilder struct {
	encoder   *sliceEncoder
	elemRType rtype
}

func newSliceEncoder(typ reflect.Type) *sliceEncoderBuilder {
	elemType := typ.Elem()
	return &sliceEncoderBuilder{
		encoder: &sliceEncoder{
			elemSize: elemType.Size(),
		},
		elemRType: rtypeOfType(elemType),
	}
}

type sliceEncoder struct {
	elemSize    uintptr
	elemEncoder ValEncoder
}

func (enc *sliceEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if s.Error != nil {
		return
	}
	if ptr == nil {
		s.Null()
		return
	}
	sh := (*reflect.SliceHeader)(ptr)
	if sh.Data == 0 {
		s.Null()
		return
	}
	if sh.Len == 0 {
		s.RawString("[]")
		return
	}
	s.ArrayStart()
	curPtr := sh.Data
	for i := 0; i < sh.Len; i++ {
		enc.elemEncoder.Encode(unsafe.Pointer(curPtr), s)
		if s.Error != nil {
			return
		}
		curPtr += enc.elemSize
	}
	s.ArrayEnd()
}
