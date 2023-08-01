package jzon

import (
	"encoding"
)

type textMarshalerValEncoder struct{}

func (textMarshalerValEncoder) Encode(o interface{}, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
	b, err := o.(encoding.TextMarshaler).MarshalText()
	if err != nil {
		s.Error = err
		return
	}
	s.String(localByteToString(b))
}
