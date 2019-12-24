package jzon

import (
	"encoding"
	"reflect"
	"unsafe"
)

var (
	textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

type textMarshalerEncoder rtype

func (enc textMarshalerEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if ptr == nil {
		s.Null()
		return
	}
	obj := packEFace(rtype(enc), ptr)
	marshaler := obj.(encoding.TextMarshaler)
	b, err := marshaler.MarshalText()
	if err != nil {
		s.Error = err
		return
	}
	s.String(localByteToString(b))
}

type dynamicTextMarshalerEncoder struct{}

func (*dynamicTextMarshalerEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	marshaler := *(*encoding.TextMarshaler)(ptr)
	if marshaler == nil {
		s.Null()
		return
	}
	b, err := marshaler.MarshalText()
	if err != nil {
		s.Error = err
		return
	}
	s.String(localByteToString(b))
}
