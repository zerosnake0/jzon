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

func (enc *pointerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *(*unsafe.Pointer)(ptr) == nil
}

func (enc *pointerEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	enc.encoder.Encode(*(*unsafe.Pointer)(ptr), s, opts)
}
