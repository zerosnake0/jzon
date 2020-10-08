package jzon

import (
	"io"
	"runtime/debug"
	"testing"
)

func TestValDecoder_Native_Ptr(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		t.Log(">>>>> initValues >>>>>")
		printValue(t, "p1", p1)
		printValue(t, "p2", p2)
		t.Log(">>>>>>>>>>>>>>>>>>>>>>")
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
		t.Log("<<<<< initValues <<<<<")
		printValue(t, "p1", p1)
		printValue(t, "p2", p2)
		t.Log("<<<<<<<<<<<<<<<<<<<<<<")
	}
	t.Run("nil pointer", func(t *testing.T) {
		f(t, "null", ErrNilPointerReceiver, (*int)(nil), (*int)(nil))
	})
	t.Run("eof", func(t *testing.T) {
		pi1 := (*int)(nil)
		pi2 := (*int)(nil)
		f(t, "", io.EOF, &pi1, &pi2)
	})
	t.Run("invalid null", func(t *testing.T) {
		pi1 := (*int)(nil)
		pi2 := (*int)(nil)
		f(t, "nul", io.EOF, &pi1, &pi2)
	})
	t.Run("null 1", func(t *testing.T) {
		pi1 := (*int)(nil)
		pi2 := (*int)(nil)
		f(t, "null", nil, &pi1, &pi2)
	})
	t.Run("null 2", func(t *testing.T) {
		i1 := 1
		i2 := 1
		pi1 := &i1
		pi2 := &i2
		f(t, "null", nil, &pi1, &pi2)
	})
	t.Run("not null error 1", func(t *testing.T) {
		pi1 := (*int)(nil)
		pi2 := (*int)(nil)
		f(t, `true`, InvalidDigitError{}, &pi1, &pi2)
	})
	t.Run("not null error 2", func(t *testing.T) {
		i1 := 1
		i2 := 1
		pi1 := &i1
		pi2 := &i2
		f(t, `true`, InvalidDigitError{}, &pi1, &pi2)
	})
	t.Run("not null 1", func(t *testing.T) {
		pi1 := (*int)(nil)
		pi2 := (*int)(nil)
		f(t, `23`, nil, &pi1, &pi2)
	})
	t.Run("not null 2", func(t *testing.T) {
		i1 := 1
		i2 := 1
		pi1 := &i1
		pi2 := &i2
		f(t, `23`, nil, &pi1, &pi2)
	})
	debug.FreeOSMemory()
}
