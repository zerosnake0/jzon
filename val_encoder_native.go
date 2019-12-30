package jzon

import (
	"unsafe"
)

// bool encoder
type boolEncoder struct{}

func (*boolEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if ptr == nil {
		s.Null()
		return
	}
	s.Bool(*(*bool)(ptr))
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
