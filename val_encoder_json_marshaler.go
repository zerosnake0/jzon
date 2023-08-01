package jzon

import (
	"encoding/json"
)

type jsonMarshalerValEncoder struct{}

func (jsonMarshalerValEncoder) Encode(o interface{}, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
	raw, err := o.(json.Marshaler).MarshalJSON()
	if err != nil {
		s.Error = err
		return
	}
	s.Raw(raw)
}
