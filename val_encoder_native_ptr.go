package jzon

import (
	"log"
	"reflect"
	"unsafe"
)

type pointerEncoderBuilder struct {
	encoder *pointerEncoder
	rtype   rtype
}

func newPointerEncoder(elemType reflect.Type) *pointerEncoderBuilder {
	return &pointerEncoderBuilder{
		encoder: &pointerEncoder{},
		rtype:   rtypeOfType(elemType),
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
	log.Printf("ptr %x", ptr)
	enc.encoder.Encode(*(*unsafe.Pointer)(ptr), s)
}
