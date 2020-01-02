package jzon

import (
	"encoding/json"
	"testing"
)

func TestValEncoder_JsonNumber(t *testing.T) {
	t.Run("non pointer", func(t *testing.T) {
		f := func(t *testing.T, n json.Number) {
			checkEncodeValueWithStandard(t, n, nil)
		}
		t.Run("empty", func(t *testing.T) {
			f(t, "")
		})
		t.Run("non empty", func(t *testing.T) {
			f(t, "-1.2e-3")
		})
		t.Run("invalid", func(t *testing.T) {
			// TODO:
		})
	})
	t.Run("pointer", func(t *testing.T) {
		f := func(t *testing.T, ptr *json.Number, err error) {
			checkEncodeValueWithStandard(t, ptr, err)
		}
		t.Run("nil", func(t *testing.T) {
			f(t, nil, nil)
		})
		t.Run("non nil", func(t *testing.T) {
			var n json.Number = "1.23"
			f(t, &n, nil)
		})
	})
	t.Run("pointer of pointer", func(t *testing.T) {
		f := func(t *testing.T, ptr **json.Number, err error) {
			checkEncodeValueWithStandard(t, ptr, err)
		}
		t.Run("nil", func(t *testing.T) {
			f(t, nil, nil)
		})
		t.Run("pointer of nil", func(t *testing.T) {
			ptr := (*json.Number)(nil)
			f(t, &ptr, nil)
		})
		t.Run("pointer of non nil", func(t *testing.T) {
			var n json.Number = "1.23"
			ptr := &n
			f(t, &ptr, nil)
		})
	})
}
