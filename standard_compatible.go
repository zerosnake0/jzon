package jzon

import (
	"io"
)

func Unmarshal(data []byte, o interface{}) error {
	return DefaultDecoder.Unmarshal(data, o)
}

func Marshal(o interface{}) ([]byte, error) {
	return DefaultEncoder.Marshal(o)
}

func Valid(data []byte) bool {
	it := NewIterator()
	b := it.Valid(data)
	ReturnIterator(it)
	return b
}

func UnmarshalFromReader(r io.Reader, o interface{}) error {
	return DefaultDecoder.UnmarshalFromReader(r, o)
}
