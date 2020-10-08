package jzon

import (
	"errors"
	"io"
	"testing"
)

type testJSONUnmarshaler struct {
	data string
	err  error
}

func (t *testJSONUnmarshaler) UnmarshalJSON(data []byte) error {
	t.data = string(data)
	return t.err
}

func TestValDecoder_JsonUnmarshaler(t *testing.T) {
	f := func(t *testing.T, data string, ex error, initValue string, initErr error) {
		var p1 *testJSONUnmarshaler
		var p2 *testJSONUnmarshaler
		if initValue != "" {
			b1 := testJSONUnmarshaler{
				data: initValue,
				err:  initErr,
			}
			p1 = &b1
			b2 := testJSONUnmarshaler{
				data: initValue,
				err:  initErr,
			}
			p2 = &b2
		}
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
	}
	f2 := func(t *testing.T, data string, ex error, initErr error) {
		f(t, data, ex, "dummy", initErr)
	}
	t.Run("nil pointer", func(t *testing.T) {
		f(t, "null", ErrNilPointerReceiver, "", nil)
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
