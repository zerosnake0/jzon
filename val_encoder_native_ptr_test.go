package jzon

import "testing"

func TestValEncoder_Ptr(t *testing.T) {
	f := func(t *testing.T, b *bool) {
		checkEncodeValueWithStandard(t, DefaultEncoder, b)
	}
	t.Run("nil", func(t *testing.T) {
		f(t, nil)
	})
	t.Run("true", func(t *testing.T) {
		b := true
		f(t, &b)
	})
	t.Run("false", func(t *testing.T) {
		b := false
		f(t, &b)
	})
}
