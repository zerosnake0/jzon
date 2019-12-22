package jzon

import (
	"encoding/json"
	"testing"
)

func TestValEncoder_JsonRawMessage(t *testing.T) {
	f := func(t *testing.T, s string) {
		msg := json.RawMessage(s)
		checkEncodeValueWithStandard(t, DefaultEncoder, msg)
	}
	t.Run("null", func(t *testing.T) {
		f(t, "null")
	})
}
