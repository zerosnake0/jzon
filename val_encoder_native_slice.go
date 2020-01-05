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

func (enc *sliceEncoder) Encode(ptr unsafe.Pointer, s *Streamer, _ *EncOpts) {
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
		enc.elemEncoder.Encode(unsafe.Pointer(curPtr), s, nil)
		if s.Error != nil {
			return
		}
		curPtr += enc.elemSize
	}
	s.ArrayEnd()
}

type sliceEncoderBuilder2 struct {
	encoder  *sliceEncoder2
	elemType reflect.Type
}

func newSliceEncoder2(typ reflect.Type) *sliceEncoderBuilder2 {
	elemType := typ.Elem()
	return &sliceEncoderBuilder2{
		encoder:  &sliceEncoder2{},
		elemType: elemType,
	}
}

type sliceEncoder2 struct {
	elemEncoder ValEncoder2
}

func (enc *sliceEncoder2) Encode2(v reflect.Value, s *Streamer, _ *EncOpts) {
	if s.Error != nil {
		return
	}
	if v.IsNil() {
		s.Null()
		return
	}
	s.ArrayStart()
	l := v.Len()
	i := 0
	for {
		enc.elemEncoder.Encode2(v.Index(i), s, nil)
		if s.Error != nil {
			return
		}
		i++
		if i == l {
			break
		}
	}
	s.ArrayEnd()
}
