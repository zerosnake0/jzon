package jzon

import (
	"bytes"
	"io"
)

// Decoder is almost standard library compatible
// The following standard methods are not implemented
// - Token
type Decoder struct {
	it  *Iterator
	err error
}

// Release decoder, decoder should not be reused after call
func (dec *Decoder) Release() {
	dec.it.Release()
	dec.it = nil
}

// UseNumber causes the Decoder to unmarshal a number into an interface{} as a
// Number instead of as a float64.
func (dec *Decoder) UseNumber() {
	dec.it.useNumber = true
}

// DisallowUnknownFields causes the Decoder to return an error when the destination
// is a struct and the input contains object keys which do not match any
// non-ignored, exported fields in the destination.
func (dec *Decoder) DisallowUnknownFields() {
	dec.it.disallowUnknownFields = true
}

// Decode reads the next JSON-encoded value from its
// input and stores it in the value pointed to by v.
func (dec *Decoder) Decode(v interface{}) error {
	if dec.err == nil {
		dec.err = dec.it.ReadVal(v)
	}
	return dec.err
}

// Buffered returns a reader of the data remaining in the Decoder's
// buffer. The reader is valid until the next call to Decode.
func (dec *Decoder) Buffered() io.Reader {
	return bytes.NewReader(dec.it.Buffer())
}

// func (dec *Decoder) Token() (json.Token, error) {
// 	panic("not implemented")
// }

// More reports whether there is another element in the
// current array or object being parsed.
func (dec *Decoder) More() bool {
	if dec.err != nil {
		return false
	}
	_, err := dec.it.nextToken()
	return err == nil
}

// InputOffset returns the input stream byte offset of the current decoder position.
// The offset gives the location of the end of the most recently returned token
// and the beginning of the next token.
// Whitespace may present at the position of offset
func (dec *Decoder) InputOffset() int64 {
	return int64(dec.it.offset + dec.it.head)
}
