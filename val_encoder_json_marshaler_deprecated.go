package jzon

import (
	"encoding/json"
	"reflect"
	"unsafe"
)

var (
	jsonMarshalerType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
)

// Deprecated
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

// Deprecated
type directJSONMarshalerEncoder struct {
	isEmpty isEmptyFunc
	rtype   rtype
}

func (enc *directJSONMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return enc.isEmpty(ptr)
}

func (enc *directJSONMarshalerEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
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

// Deprecated
type pointerJSONMarshalerEncoder rtype

func (enc pointerJSONMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *(*unsafe.Pointer)(ptr) == nil
}

func (enc pointerJSONMarshalerEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
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

// Deprecated
type dynamicJSONMarshalerEncoder struct{}

func (*dynamicJSONMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *(*json.Marshaler)(ptr) == nil
}

func (*dynamicJSONMarshalerEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
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
