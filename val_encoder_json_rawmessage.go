package jzon

import (
	"unsafe"
)

type jsonRawMessageEncoder struct {
}

func (*jsonRawMessageEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	data := *(*[]byte)(ptr)
	s.Raw(data)
}
