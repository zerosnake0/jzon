package jzon

import (
	"encoding/json"
	"io"
	"testing"
)

func TestValDecoder_JsonRawMessage(t *testing.T) {
	f := func(t *testing.T, data string, ex error, initValue string) {
		var p1 *json.RawMessage
		var p2 *json.RawMessage
		if initValue != "" {
			b1 := append(json.RawMessage(nil), initValue...)
			p1 = &b1
			b2 := append(json.RawMessage(nil), initValue...)
			p2 = &b2
		}
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
	}
	f2 := func(t *testing.T, data string, ex error) {
		f(t, data, ex, "1.23")
	}
	t.Run("nil pointer", func(t *testing.T) {
		f(t, "null", ErrNilPointerReceiver, "")
	})
	t.Run("eof", func(t *testing.T) {
		f2(t, "", io.EOF)
	})
	t.Run("valid", func(t *testing.T) {
		f2(t, `null`, nil)
	})
}
