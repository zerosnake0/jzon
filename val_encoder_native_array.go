package jzon

import (
	"log"
	"reflect"
	"unsafe"
)

type arrayEncoderBuilder struct {
	encoder   *arrayEncoder
	elemRType rtype
}

func newArrayEncoder(typ reflect.Type) *arrayEncoderBuilder {
	elemType := typ.Elem()
	return &arrayEncoderBuilder{
		encoder: &arrayEncoder{
			elemSize: elemType.Size(),
			length:   typ.Len(),
		},
		elemRType: rtypeOfType(elemType),
	}
}

type arrayEncoder struct {
	encoder  ValEncoder
	elemSize uintptr
	length   int
}

func (enc *arrayEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if ptr == nil {
		s.Null()
		return
	}
	log.Println(ptr)
	s.ArrayStart()
	for i := 0; i < enc.length; i++ {
		enc.encoder.Encode(ptr, s)
		if s.Error != nil {
			return
		}
		ptr = add(ptr, enc.elemSize, "i < enc.length")
	}
	s.ArrayEnd()
}
