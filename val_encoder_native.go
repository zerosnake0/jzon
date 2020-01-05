package jzon

import (
	"reflect"
	"unsafe"
)

// bool encoder
type boolEncoder struct{}

func (*boolEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	quoted := (opts != nil) && opts.Quoted
	if !quoted {
		s.Bool(*(*bool)(ptr))
		return
	}
	if *(*bool)(ptr) {
		s.String("true")
	} else {
		s.String("false")
	}
}

func (*boolEncoder) Encode2(v reflect.Value, s *Streamer, opts *EncOpts) {
	quoted := (opts != nil) && opts.Quoted
	b := v.Bool()
	if !quoted {
		s.Bool(b)
		return
	}
	if b {
		s.String("true")
	} else {
		s.String("false")
	}
}

// string encoder
type stringEncoder struct{}

func (*stringEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	quoted := (opts != nil) && opts.Quoted
	if !quoted {
		s.String(*(*string)(ptr))
		return
	}
	subStream := s.encoder.NewStreamer()
	defer s.encoder.ReturnStreamer(subStream)
	subStream.String(*(*string)(ptr))
	if subStream.Error != nil {
		s.Error = subStream.Error
		return
	}
	s.String(localByteToString(subStream.buffer))
}

func (*stringEncoder) Encode2(v reflect.Value, s *Streamer, opts *EncOpts) {
	quoted := (opts != nil) && opts.Quoted
	str := v.String()
	if !quoted {
		s.String(str)
		return
	}
	subStream := s.encoder.NewStreamer()
	defer s.encoder.ReturnStreamer(subStream)
	subStream.String(str)
	if subStream.Error != nil {
		s.Error = subStream.Error
		return
	}
	s.String(localByteToString(subStream.buffer))
}
