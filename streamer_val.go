package jzon

import (
	"reflect"
	"unsafe"
)

func (s *Streamer) Value(obj interface{}) *Streamer {
	if s.Error != nil {
		return s
	}
	if obj == nil {
		s.Null()
		return s
	}
	eface := (*eface)(unsafe.Pointer(&obj))
	enc := s.encoder.getEncoderFromCache(eface.rtype)
	if enc == nil {
		typ := reflect.TypeOf(obj)
		enc = s.encoder.createEncoder(eface.rtype, typ)
	}
	enc.Encode(eface.data, s, nil)
	return s
}
