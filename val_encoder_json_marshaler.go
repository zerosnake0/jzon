package jzon

import (
	"encoding/json"
	"reflect"
	"unsafe"
)

var (
	jsonMarshalerType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
)

type jsonMarshalerEncoder rtype

func (enc jsonMarshalerEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
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

type directJsonMarshalerEncoder rtype

func (enc directJsonMarshalerEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
	if ptr == nil {
		s.Null()
		return
	}
	obj := packEFace(rtype(enc), *(*unsafe.Pointer)(ptr))
	marshaler := obj.(json.Marshaler)
	b, err := marshaler.MarshalJSON()
	if err != nil {
		s.Error = err
		return
	}
	s.Raw(b)
}

type dynamicJsonMarshalerEncoder struct{}

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
