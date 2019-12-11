package jzon

import (
	"errors"
	"io"
	"testing"
)

type testJsonUnmarshaler struct {
	data string
	err  error
}

func (t *testJsonUnmarshaler) UnmarshalJSON(data []byte) error {
	t.data = string(data)
	return t.err
}

func TestValDecoder_JsonUnmarshaler(t *testing.T) {
	f := func(t *testing.T, data string, ex error, initValue string, initErr error) {
		var p1 *testJsonUnmarshaler
		var p2 *testJsonUnmarshaler
		if initValue != "" {
			b1 := testJsonUnmarshaler{
				data: initValue,
				err:  initErr,
			}
			p1 = &b1
			b2 := testJsonUnmarshaler{
				data: initValue,
				err:  initErr,
			}
			p2 = &b2
		}
		checkStandard(t, DefaultDecoder, data, ex, p1, p2)
	}
	f2 := func(t *testing.T, data string, ex error, initErr error) {
		f(t, data, ex, "dummy", initErr)
	}
	t.Run("nil pointer", func(t *testing.T) {
		f(t, "null", NilPointerReceiverError, "", nil)
	})
	t.Run("eof", func(t *testing.T) {
		f2(t, " ", io.EOF, nil)
	})
	t.Run("custom error", func(t *testing.T) {
		f2(t, ` null `, errors.New("test"), errors.New("test"))
	})
	t.Run("no error", func(t *testing.T) {
		f2(t, ` null `, nil, nil)
	})
}
