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

func (enc textMarshalerEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
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

type directTextMarshalerEncoder rtype

func (enc directTextMarshalerEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
	if ptr == nil {
		s.Null()
		return
	}
	obj := packEFace(rtype(enc), *(*unsafe.Pointer)(ptr))
	marshaler := obj.(encoding.TextMarshaler)
	b, err := marshaler.MarshalText()
	if err != nil {
		s.Error = err
		return
	}
	s.String(localByteToString(b))
}

type dynamicTextMarshalerEncoder struct{}

func (*dynamicTextMarshalerEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
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
