package jzon

import (
	"reflect"
	"unsafe"
)

type pointerEncoderBuilder struct {
	encoder   *pointerEncoder
	elemRType rtype
}

func newPointerEncoder(elemType reflect.Type) *pointerEncoderBuilder {
	return &pointerEncoderBuilder{
		encoder:   &pointerEncoder{},
		elemRType: rtypeOfType(elemType),
	}
}

type pointerEncoder struct {
	encoder ValEncoder
}

func (enc *pointerEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	enc.encoder.Encode(*(*unsafe.Pointer)(ptr), s, opts)
}

type pointerEncoderBuilder2 struct {
	encoder  *pointerEncoder2
	elemType reflect.Type
}

func newPointerEncoder2(elemType reflect.Type) *pointerEncoderBuilder2 {
	return &pointerEncoderBuilder2{
		encoder:  &pointerEncoder2{},
		elemType: elemType,
	}
}

type pointerEncoder2 struct {
	encoder ValEncoder2
}

func (enc *pointerEncoder2) Encode2(v reflect.Value, s *Streamer, opts *EncOpts) {
	ptr := ptrOfValue(v)
	if ptr == nil {
		s.Null()
		return
	}
	enc.encoder.Encode2(v.Elem(), s, opts)
}
