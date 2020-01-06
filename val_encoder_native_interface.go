package jzon

import (
	"reflect"
	"unsafe"
)

type efaceEncoder struct{}

func (*efaceEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	s.Value(*(*interface{})(ptr))
}

func (*efaceEncoder) Encode2(v reflect.Value, s *Streamer, opts *EncOpts) {
	if v.IsNil() {
		s.Null()
		return
	}
	s.value2(v.Elem().Interface(), opts)
}

type ifaceEncoder struct{}

func (*ifaceEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	o := packIFace(ptr)
	s.Value(o)
}
