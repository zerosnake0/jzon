package jzon

import (
	"reflect"
	"unsafe"
)

type sliceEncoderBuilder struct {
	encoder  *sliceEncoder
	elemType reflect.Type
}

func newSliceEncoder(typ reflect.Type) *sliceEncoderBuilder {
	elemType := typ.Elem()
	return &sliceEncoderBuilder{
		encoder: &sliceEncoder{
			elemSize: elemType.Size(),
		},
		elemType: elemType,
	}
}

type sliceEncoder struct {
	elemSize    uintptr
	elemEncoder ValEncoder
}

func (*sliceEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	sh := (*sliceHeader)(ptr)
	return sh.Data == 0 || sh.Len == 0
}

func (enc *sliceEncoder) Encode(ptr unsafe.Pointer, s *Streamer, _ *EncOpts) {
	if s.Error != nil {
		return
	}
	if ptr == nil {
		s.Null()
		return
	}
	sh := (*sliceHeader)(ptr)
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
		enc.elemEncoder.Encode(unsafe.Pointer(curPtr), s, nil)
		if s.Error != nil {
			return
		}
		curPtr += enc.elemSize
	}
	s.ArrayEnd()
}
