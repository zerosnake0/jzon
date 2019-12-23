package jzon

import (
	"reflect"
	"unsafe"
)

type pointerEncoderBuilder struct {
	encoder   *pointerEncoder
	ptrRType  rtype
	elemRType rtype
}

func newPointerEncoder(ptrRType rtype, elemType reflect.Type) *pointerEncoderBuilder {
	return &pointerEncoderBuilder{
		encoder:   &pointerEncoder{},
		ptrRType:  ptrRType,
		elemRType: rtypeOfType(elemType),
	}
}

type pointerEncoder struct {
	encoder ValEncoder
}

func (enc *pointerEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if ptr == nil {
		s.Null()
		return
	}
	enc.encoder.Encode(*(*unsafe.Pointer)(ptr), s)
}
