package jzon

import (
	"reflect"
	"unsafe"
)

type jsonNumberEncoder struct {
}

func (*jsonNumberEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	str := *(*string)(ptr)
	if str == "" {
		str = "0"
	}
	// TODO: the standard lib will check the validity
	s.RawString(str)
}

func (*jsonNumberEncoder) Encode2(v reflect.Value, s *Streamer, opts *EncOpts) {
	ptr := ptrOfValue(v)
	str := *(*string)(ptr)
	if str == "" {
		str = "0"
	}
	s.RawString(str)
}
