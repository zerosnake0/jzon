package jzon

import (
	"io"
	"testing"
)

func TestValDecoder_Native_Float32(t *testing.T) {
	f := func(t *testing.T, data string, ex error, initValue float32) {
		var p1 *float32
		var p2 *float32
		if initValue != 0 {
			b1 := initValue
			p1 = &b1
			b2 := initValue
			p2 = &b2
		}
		checkStandard(t, DefaultDecoder, data, ex, p1, p2)
	}
	f2 := func(t *testing.T, data string, ex error) {
		f(t, data, ex, 1.234)
	}
	t.Run("nil pointer", func(t *testing.T) {
		f(t, "null", NilPointerReceiverError, 0)
	})
	t.Run("eof", func(t *testing.T) {
		f2(t, "", io.EOF)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		f2(t, `true`, InvalidFloatError{})
	})
	t.Run("invalid null", func(t *testing.T) {
		f2(t, "nul", io.EOF)
	})
	t.Run("null", func(t *testing.T) {
		f2(t, "null", nil)
	})
	t.Run("valid", func(t *testing.T) {
		f2(t, "0.123e-4", nil)
	})
}

func TestValDecoder_Native_Float64(t *testing.T) {
	f := func(t *testing.T, data string, ex error, initValue float64) {
		var p1 *float64
		var p2 *float64
		if initValue != 0 {
			b1 := initValue
			p1 = &b1
			b2 := initValue
			p2 = &b2
		}
		checkStandard(t, DefaultDecoder, data, ex, p1, p2)
	}
	f2 := func(t *testing.T, data string, ex error) {
		f(t, data, ex, 1.234)
	}
	t.Run("nil pointer", func(t *testing.T) {
		f(t, "null", NilPointerReceiverError, 0)
	})
	t.Run("eof", func(t *testing.T) {
		f2(t, "", io.EOF)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		f2(t, `true`, InvalidFloatError{})
	})
	t.Run("invalid null", func(t *testing.T) {
		f2(t, "nul", io.EOF)
	})
	t.Run("null", func(t *testing.T) {
		f2(t, "null", nil)
	})
	t.Run("valid", func(t *testing.T) {
		f2(t, "0.123e-45", nil)
	})
}
