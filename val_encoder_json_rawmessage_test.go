package jzon

import (
	"encoding/json"
	"testing"
)

func TestValEncoder_JsonRawMessage(t *testing.T) {
	t.Run("non pointer", func(t *testing.T) {
		f := func(t *testing.T, s string) {
			msg := json.RawMessage(s)
			checkEncodeValueWithStandard(t, msg, nil)
		}
		t.Run("null", func(t *testing.T) {
			f(t, "null")
		})
		t.Run("true", func(t *testing.T) {
			f(t, "true")
		})
	})
	t.Run("pointer", func(t *testing.T) {
		f := func(t *testing.T, msg *json.RawMessage, err error) {
			checkEncodeValueWithStandard(t, msg, err)
		}
		t.Run("nil", func(t *testing.T) {
			f(t, nil, nil)
		})
		t.Run("non nil", func(t *testing.T) {
			msg := json.RawMessage("false")
			f(t, &msg, nil)
		})
	})
	t.Run("pointer of pointer", func(t *testing.T) {
		f := func(t *testing.T, msg **json.RawMessage, err error) {
			checkEncodeValueWithStandard(t, msg, err)
		}
		t.Run("nil", func(t *testing.T) {
			f(t, nil, nil)
		})
		t.Run("pointer of nil", func(t *testing.T) {
			ptr := (*json.RawMessage)(nil)
			f(t, &ptr, nil)
		})
		t.Run("non nil", func(t *testing.T) {
			msg := json.RawMessage("false")
			ptr := &msg
			f(t, &ptr, nil)
		})
	})
}
