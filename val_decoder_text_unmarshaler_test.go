package jzon

import (
	"errors"
	"io"
	"testing"
)

type testTextUnmarshaler struct {
	data string
	err  error
}

func (t *testTextUnmarshaler) UnmarshalText(data []byte) error {
	t.data = string(data)
	return t.err
}

func TestValDecoder_TextUnmarshaler(t *testing.T) {
	f := func(t *testing.T, data string, ex error, initValue string, initErr error) {
		var p1 *testTextUnmarshaler
		var p2 *testTextUnmarshaler
		if initValue != "" {
			b1 := testTextUnmarshaler{
				data: initValue,
				err:  initErr,
			}
			p1 = &b1
			b2 := testTextUnmarshaler{
				data: initValue,
				err:  initErr,
			}
			p2 = &b2
		}
		checkDecodeWithStandard(t, DefaultDecoder, data, ex, p1, p2)
	}
	f2 := func(t *testing.T, data string, ex, initErr error) {
		f(t, data, ex, "dummy", initErr)
	}
	t.Run("nil pointer", func(t *testing.T) {
		f(t, "null", NilPointerReceiverError, "", nil)
	})
	t.Run("eof", func(t *testing.T) {
		f2(t, " ", io.EOF, nil)
	})
	t.Run("invalid byte", func(t *testing.T) {
		f2(t, ` + `, UnexpectedByteError{}, nil)
	})
	t.Run("null", func(t *testing.T) {
		f2(t, ` null `, nil, nil)
	})
	t.Run("invalid string", func(t *testing.T) {
		f2(t, ` " `, io.EOF, nil)
	})
	t.Run("no error", func(t *testing.T) {
		f2(t, ` "abc" `, nil, nil)
	})
	t.Run("custom error", func(t *testing.T) {
		f2(t, ` "abc" `, errors.New("test"), errors.New("test"))
	})
}
