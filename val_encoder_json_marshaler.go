package jzon

import (
	"encoding/json"
	"reflect"
	"unsafe"
)

var (
	jsonMarshalerType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
)

type jsonMarshalerEncoder struct {
	isEmpty isEmptyFunc
	rtype   rtype
}

func (enc *jsonMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return enc.isEmpty(ptr)
}

func (enc *jsonMarshalerEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
	if ptr == nil {
		s.Null()
		return
	}
	obj := packEFace(enc.rtype, ptr)
	marshaler := obj.(json.Marshaler)
	b, err := marshaler.MarshalJSON()
	if err != nil {
		s.Error = err
		return
	}
	s.Raw(b)
}

type directJsonMarshalerEncoder struct {
	isEmpty isEmptyFunc
	rtype   rtype
}

func (enc *directJsonMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return enc.isEmpty(ptr)
}

func (enc *directJsonMarshalerEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
	if ptr == nil {
		s.Null()
		return
	}
	obj := packEFace(enc.rtype, *(*unsafe.Pointer)(ptr))
	marshaler := obj.(json.Marshaler)
	b, err := marshaler.MarshalJSON()
	if err != nil {
		s.Error = err
		return
	}
	s.Raw(b)
}

type pointerJsonMarshalerEncoder rtype

func (enc pointerJsonMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *(*unsafe.Pointer)(ptr) == nil
}

func (enc pointerJsonMarshalerEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
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
	marshaler := obj.(json.Marshaler)
	b, err := marshaler.MarshalJSON()
	if err != nil {
		s.Error = err
		return
	}
	s.Raw(b)
}

type dynamicJsonMarshalerEncoder struct{}

func (*dynamicJsonMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *(*json.Marshaler)(ptr) == nil
}

func (*dynamicJsonMarshalerEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
	if ptr == nil {
		s.Null()
		return
	}
	marshaler := *(*json.Marshaler)(ptr)
	if marshaler == nil {
		s.Null()
		return
	}
	b, err := marshaler.MarshalJSON()
	if err != nil {
		s.Error = err
		return
	}
	s.Raw(b)
}
