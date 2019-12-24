package jzon

import (
	"unsafe"
)

type directEncoder struct {
	encoder ValEncoder
}

func (enc *directEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	enc.encoder.Encode(unsafe.Pointer(&ptr), s)
}