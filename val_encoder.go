package jzon

import (
	"unsafe"
)

type ValEncoder interface {
	Encode(ptr unsafe.Pointer, s *Streamer) error
}
