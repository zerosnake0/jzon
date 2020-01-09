package jzon

import (
	"encoding/json"
	"io"
	"testing"
)

func TestValDecoder_JsonNumber(t *testing.T) {
	f := func(t *testing.T, data string, ex error, initValue json.Number) {
		var p1 *json.Number
		var p2 *json.Number
		if initValue != "" {
			b1 := initValue
			p1 = &b1
			b2 := initValue
			p2 = &b2
		}
		checkDecodeWithStandard(t, DefaultDecoder, data, ex, p1, p2)
	}
	f2 := func(t *testing.T, data string, ex error) {
		f(t, data, ex, "1.23")
	}
	t.Run("nil pointer", func(t *testing.T) {
		f(t, "null", NilPointerReceiverError, "")
	})
	t.Run("eof", func(t *testing.T) {
		f2(t, "", io.EOF)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		f2(t, `+`, UnexpectedByteError{})
	})
	t.Run("invalid null", func(t *testing.T) {
		f2(t, "nul", io.EOF)
	})
	t.Run("null", func(t *testing.T) {
		f2(t, "null", nil)
	})
	t.Run("invalid string", func(t *testing.T) {
		f2(t, `"abc`, io.EOF)
	})
	t.Run("string", func(t *testing.T) {
		// TODO: this will be invalid in the future version of golang?
		f2(t, `"abc"`, nil)
	})
	t.Run("invalid number", func(t *testing.T) {
		f2(t, `-0.e`, InvalidDigitError{})
	})
	t.Run("valid number", func(t *testing.T) {
		f2(t, `123.456e+0789`, nil)
	})
}
