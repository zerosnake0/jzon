package jzon

import (
	"reflect"
	"unsafe"
)

func (s *Streamer) value(obj interface{}, opts *EncOpts) *Streamer {
	if s.Error != nil {
		return s
	}
	if obj == nil {
		s.Null()
		return s
	}
	ef := (*eface)(unsafe.Pointer(&obj))
	enc := s.cfg.getEncoderFromCache(ef.rtype)
	if enc == nil {
		typ := reflect.TypeOf(obj)
		enc = s.cfg.createEncoder(ef.rtype, typ)
	}
	enc.Encode(ef.data, s, opts)
	return s
}

// Value encodes a golang object to json
func (s *Streamer) Value(obj interface{}) *Streamer {
	return s.value(obj, nil)
}
