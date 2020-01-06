package jzon

import (
	"log"
	"reflect"
	"unsafe"
)

// special array (empty)
type emptyArrayEncoder struct{}

func (enc *emptyArrayEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	s.RawString("[]")
}

func (enc *emptyArrayEncoder) Encode2(v reflect.Value, s *Streamer, opts *EncOpts) {
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

type arrayEncoderBuilder2 struct {
	encoder  *arrayEncoder2
	elemType reflect.Type
}

func newArrayEncoder2(typ reflect.Type) *arrayEncoderBuilder2 {
	elemType := typ.Elem()
	return &arrayEncoderBuilder2{
		encoder: &arrayEncoder2{
			length: typ.Len(),
		},
		elemType: elemType,
	}
}

type arrayEncoder2 struct {
	encoder ValEncoder2
	length  int
}

func (enc *arrayEncoder2) Encode2(v reflect.Value, s *Streamer, _ *EncOpts) {
	s.ArrayStart()
	i := 0
	for {
		log.Println(v.Index(i).CanAddr())
		enc.encoder.Encode2(v.Index(i), s, nil)
		if s.Error != nil {
			return
		}
		i++
		if i == enc.length {
			break
		}
	}
	s.ArrayEnd()
}
