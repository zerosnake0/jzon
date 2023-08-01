//go:build go1.18
// +build go1.18

package jzon

import (
	"unsafe"
)

func (enc *directMapEncoder) Encode(ptr unsafe.Pointer, s *Streamer, _ *EncOpts) {
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
	s.ObjectStart()
	var iter hiter
	mapiterinit(enc.mapRType, ptr, &iter)
	for i := 0; iter.key != nil; i++ {
		enc.keyEncoder.Encode(iter.key, s)
		if s.Error != nil {
			return
		}
		enc.elemEncoder.Encode(iter.elem, s, nil)
		if s.Error != nil {
			return
		}
		mapiternext(&iter)
	}
	s.ObjectEnd()
}
