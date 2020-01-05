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
	ef := (*eface)(unsafe.Pointer(&obj))
	enc := s.encoder.getEncoderFromCache(ef.rtype)
	if enc == nil {
		typ := reflect.TypeOf(obj)
		enc = s.encoder.createEncoder(ef.rtype, typ)
	}
	enc.Encode(ef.data, s, nil)
	return s
}

func (s *Streamer) Value2(obj interface{}) *Streamer {
	if s.Error != nil {
		return s
	}
	if obj == nil {
		s.Null()
		return s
	}
	v := reflect.ValueOf(obj)
	typ := v.Type()
	enc := s.encoder.getEncoderFromCache2(typ)
	if enc == nil {
		enc = s.encoder.createEncoder2(typ)
	}
	enc.Encode2(v, s, nil)
	return s
}
