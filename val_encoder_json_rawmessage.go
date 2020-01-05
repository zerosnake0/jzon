package jzon

import (
	"reflect"
	"unsafe"
)

type jsonRawMessageEncoder struct {
}

func (*jsonRawMessageEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	data := *(*[]byte)(ptr)
	// TODO: raw message validation?
	s.Raw(data)
}

func (*jsonRawMessageEncoder) Encode2(v reflect.Value, s *Streamer, opts *EncOpts) {
	ptr := ptrOfValue(v)
	data := *(*[]byte)(ptr)
	s.Raw(data)
}
