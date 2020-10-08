package jzon

import (
	"io"
	"testing"
)

func TestValDecoder_Native_Bool(t *testing.T) {
	f := func(t *testing.T, data string, ex error, initValue, initBool bool) {
		var p1 *bool
		var p2 *bool
		if initValue {
			b1 := initBool
			p1 = &b1
			b2 := initBool
			p2 = &b2
		}
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
	}
	f2 := func(t *testing.T, data string, ex error) {
		f(t, data, ex, true, true)
	}
	f3 := func(t *testing.T, data string, ex error) {
		f(t, data, ex, true, false)
	}
	t.Run("nil pointer", func(t *testing.T) {
		f(t, "null", NilPointerReceiverError, false, false)
	})
	t.Run("eof", func(t *testing.T) {
		f2(t, "", io.EOF)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		f2(t, `+`, UnexpectedByteError{})
	})
	t.Run("invalid null", func(t *testing.T) {
		f2(t, "nul", io.EOF)
	})
	t.Run("null", func(t *testing.T) {
		f2(t, "null", nil)
	})
	t.Run("invalid true", func(t *testing.T) {
		t.Run("init true", func(t *testing.T) {
			f2(t, `tru`, io.EOF)
		})
		t.Run("init false", func(t *testing.T) {
			f3(t, `tru`, io.EOF)
		})
	})
	t.Run("valid true", func(t *testing.T) {
		t.Run("init true", func(t *testing.T) {
			f2(t, `true`, nil)
		})
		t.Run("init false", func(t *testing.T) {
			f3(t, `true`, nil)
		})
	})
	t.Run("invalid false", func(t *testing.T) {
		t.Run("init true", func(t *testing.T) {
			f2(t, `fals`, io.EOF)
		})
	})
	t.Run("valid false", func(t *testing.T) {
		t.Run("init true", func(t *testing.T) {
			f2(t, `false`, nil)
		})
	})
}

func TestValDecoder_Native_String(t *testing.T) {
	f := func(t *testing.T, data string, ex error, initValue string) {
		var p1 *string
		var p2 *string
		if initValue != "" {
			b1 := initValue
			p1 = &b1
			b2 := initValue
			p2 = &b2
		}
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
	}
	f2 := func(t *testing.T, data string, ex error) {
		f(t, data, ex, "dummy")
	}
	t.Run("nil pointer", func(t *testing.T) {
		f(t, "null", NilPointerReceiverError, "")
	})
	t.Run("eof", func(t *testing.T) {
		f2(t, "", io.EOF)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		f2(t, `true`, UnexpectedByteError{})
	})
	t.Run("invalid null", func(t *testing.T) {
		f2(t, "nul", io.EOF)
	})
	t.Run("null", func(t *testing.T) {
		f2(t, "null", nil)
	})
	t.Run("invalid string", func(t *testing.T) {
		f2(t, `"abc`, io.EOF)
	})
	t.Run("valid string", func(t *testing.T) {
		f2(t, `"abc"`, nil)
	})
}
