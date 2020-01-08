package jzon

import (
	"reflect"
	"unsafe"
)

// special array (empty)
type emptyArrayEncoder struct{}

func (*emptyArrayEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return true
}

func (*emptyArrayEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	s.RawString("[]")
}

type arrayEncoderBuilder struct {
	encoder   *arrayEncoder
	elemType  reflect.Type
	elemRType rtype
}

func newArrayEncoder(typ reflect.Type) *arrayEncoderBuilder {
	elemType := typ.Elem()
	return &arrayEncoderBuilder{
		encoder: &arrayEncoder{
			elemSize: elemType.Size(),
			length:   typ.Len(),
		},
		elemType:  elemType,
		elemRType: rtypeOfType(elemType),
	}
}

type arrayEncoder struct {
	encoder  ValEncoder
	elemSize uintptr
	length   int
}

func (enc *arrayEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func (enc *arrayEncoder) Encode(ptr unsafe.Pointer, s *Streamer, _ *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	s.ArrayStart()
	i := 0
	for {
		enc.encoder.Encode(ptr, s, nil)
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
