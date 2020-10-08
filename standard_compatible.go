package jzon

import (
	"io"
)

func Unmarshal(data []byte, o interface{}) error {
	return DefaultDecoderConfig.Unmarshal(data, o)
}

func Marshal(o interface{}) ([]byte, error) {
	return DefaultEncoderConfig.Marshal(o)
}

func Valid(data []byte) bool {
	it := NewIterator()
	b := it.Valid(data)
	it.Release()
	return b
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return DefaultDecoderConfig.NewDecoder(r)
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return DefaultEncoderConfig.NewEncoder(w)
}
