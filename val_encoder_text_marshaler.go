package jzon

import (
	"encoding"
	"reflect"
	"unsafe"
)

var (
	textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

type textMarshalerEncoder struct {
	isEmpty isEmptyFunc
	rtype   rtype
}

func (enc *textMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return enc.isEmpty(ptr)
}

func (enc *textMarshalerEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
	if ptr == nil {
		s.Null()
		return
	}
	obj := packEFace(enc.rtype, ptr)
	marshaler := obj.(encoding.TextMarshaler)
	b, err := marshaler.MarshalText()
	if err != nil {
		s.Error = err
		return
	}
	s.String(localByteToString(b))
}

type directTextMarshalerEncoder struct {
	isEmpty isEmptyFunc
	rtype   rtype
}

func (enc *directTextMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return enc.isEmpty(ptr)
}

func (enc *directTextMarshalerEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
	if ptr == nil {
		s.Null()
		return
	}
	obj := packEFace(enc.rtype, *(*unsafe.Pointer)(ptr))
	marshaler := obj.(encoding.TextMarshaler)
	b, err := marshaler.MarshalText()
	if err != nil {
		s.Error = err
		return
	}
	s.String(localByteToString(b))
}

type pointerTextMarshalerEncoder rtype

func (enc pointerTextMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *(*unsafe.Pointer)(ptr) == nil
}

func (enc pointerTextMarshalerEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
	if ptr == nil {
		s.Null()
		return
	}
	ptr = *(*unsafe.Pointer)(ptr)
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

func (*dynamicTextMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *(*encoding.TextMarshaler)(ptr) == nil
}

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
