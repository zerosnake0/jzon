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

func (enc jsonMarshalerEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	obj := packEFace(rtype(enc), ptr)
	marshaler := obj.(json.Marshaler)
	b, err := marshaler.MarshalJSON()
	if err != nil {
		s.Error = err
		return
	}
	s.Raw(b)
}
