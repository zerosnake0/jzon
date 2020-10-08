package jzon

import (
	"bytes"
	"io"
)

type Decoder struct {
	it  *Iterator
	err error
}

func (dec *Decoder) UseNumber() {
	dec.it.useNumber = true
}

func (dec *Decoder) DisallowUnknownFields() {
	dec.it.disallowUnknownFields = true
}

func (dec *Decoder) Decode(v interface{}) error {
	if dec.err == nil {
		dec.err = dec.it.ReadVal(v)
	}
	return dec.err
}

func (dec *Decoder) Buffered() io.Reader {
	return bytes.NewReader(dec.it.Buffer())
}

// func (dec *Decoder) Token() (json.Token, error) {
// 	panic("implement me")
// }

func (dec *Decoder) More() bool {
	if dec.err != nil {
		return false
	}
	_, err := dec.it.nextToken()
	return err == nil
}

func (dec *Decoder) InputOffset() int64 {
	return int64(dec.it.offset + dec.it.head)
}
