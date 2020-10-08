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
	ReturnIterator(it)
	return b
}

func UnmarshalFromReader(r io.Reader, o interface{}) error {
	return DefaultDecoderConfig.UnmarshalFromReader(r, o)
}

func NewDecoder(r io.Reader) *Decoder {
	return DefaultDecoderConfig.NewDecoder(r)
}

func ReturnDecoder(dec *Decoder) {
	DefaultDecoderConfig.ReturnDecoder(dec)
}
