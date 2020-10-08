package jzon

import (
	"io"
)

// Unmarshal parses the JSON-encoded data and stores the result
// in the value pointed to
func Unmarshal(data []byte, o interface{}) error {
	return DefaultDecoderConfig.Unmarshal(data, o)
}

// Marshal returns the JSON encoding of object
func Marshal(o interface{}) ([]byte, error) {
	return DefaultEncoderConfig.Marshal(o)
}

// Valid reports whether data is a valid JSON encoding.
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
