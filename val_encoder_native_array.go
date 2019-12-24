package jzon

import (
	"reflect"
	"unsafe"
)

// special array (empty)
type emptyArrayEncoder struct{}

func (enc *emptyArrayEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if ptr == nil {
		s.Null()
		return
	}
	s.RawString("[]")
}

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
	s.ArrayStart()
	i := 0
	for {
		enc.encoder.Encode(ptr, s)
		if s.Error != nil {
			return
		}
		i++
		if i == enc.length {
			break
		}
		ptr = add(ptr, enc.elemSize, "i < enc.length")
	}
	s.ArrayEnd()
}