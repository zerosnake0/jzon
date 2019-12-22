package jzon

import (
	"unsafe"
)

type jsonRawMessageEncoder struct {
}

func (*jsonRawMessageEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	data := *(*[]byte)(ptr)
	// TODO: raw message validation?
	s.Raw(data)
}
