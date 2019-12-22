package jzon

import (
	"encoding/json"
	"testing"
)

func TestValEncoder_JsonNumber(t *testing.T) {
	f := func(t *testing.T, n json.Number) {
		checkEncodeValueWithStandard(t, DefaultEncoder, n)
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
}
