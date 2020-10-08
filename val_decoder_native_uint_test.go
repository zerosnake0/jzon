package jzon

import (
	"io"
	"testing"
)

func TestValDecoder_Native_Uint(t *testing.T) {
	f := func(t *testing.T, data string, ex error, initValue uint) {
		var p1 *uint
		var p2 *uint
		if initValue != 0 {
			b1 := initValue
			p1 = &b1
			b2 := initValue
			p2 = &b2
		}
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
	}
	f2 := func(t *testing.T, data string, ex error) {
		f(t, data, ex, 1)
	}
	t.Run("nil pointer", func(t *testing.T) {
		f(t, "null", ErrNilPointerReceiver, 0)
	})
	t.Run("eof", func(t *testing.T) {
		f2(t, "", io.EOF)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		f2(t, `true`, InvalidDigitError{})
	})
	t.Run("invalid null", func(t *testing.T) {
		f2(t, "nul", io.EOF)
	})
	t.Run("null", func(t *testing.T) {
		f2(t, "null", nil)
	})
	t.Run("valid", func(t *testing.T) {
		f2(t, "127", nil)
	})
}

func TestValDecoder_Native_Uint8(t *testing.T) {
	f := func(t *testing.T, data string, ex error, initValue uint8) {
		var p1 *uint8
		var p2 *uint8
		if initValue != 0 {
			b1 := initValue
			p1 = &b1
			b2 := initValue
			p2 = &b2
		}
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
	}
	f2 := func(t *testing.T, data string, ex error) {
		f(t, data, ex, 1)
	}
	t.Run("nil pointer", func(t *testing.T) {
		f(t, "null", ErrNilPointerReceiver, 0)
	})
	t.Run("eof", func(t *testing.T) {
		f2(t, "", io.EOF)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		f2(t, `true`, InvalidDigitError{})
	})
	t.Run("invalid null", func(t *testing.T) {
		f2(t, "nul", io.EOF)
	})
	t.Run("null", func(t *testing.T) {
		f2(t, "null", nil)
	})
	t.Run("valid", func(t *testing.T) {
		f2(t, "127", nil)
	})
}

func TestValDecoder_Native_Uint16(t *testing.T) {
	f := func(t *testing.T, data string, ex error, initValue uint16) {
		var p1 *uint16
		var p2 *uint16
		if initValue != 0 {
			b1 := initValue
			p1 = &b1
			b2 := initValue
			p2 = &b2
		}
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
	}
	f2 := func(t *testing.T, data string, ex error) {
		f(t, data, ex, 1)
	}
	t.Run("nil pointer", func(t *testing.T) {
		f(t, "null", ErrNilPointerReceiver, 0)
	})
	t.Run("eof", func(t *testing.T) {
		f2(t, "", io.EOF)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		f2(t, `true`, InvalidDigitError{})
	})
	t.Run("invalid null", func(t *testing.T) {
		f2(t, "nul", io.EOF)
	})
	t.Run("null", func(t *testing.T) {
		f2(t, "null", nil)
	})
	t.Run("valid", func(t *testing.T) {
		f2(t, "127", nil)
	})
}

func TestValDecoder_Native_Uint32(t *testing.T) {
	f := func(t *testing.T, data string, ex error, initValue uint32) {
		var p1 *uint32
		var p2 *uint32
		if initValue != 0 {
			b1 := initValue
			p1 = &b1
			b2 := initValue
			p2 = &b2
		}
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
	}
	f2 := func(t *testing.T, data string, ex error) {
		f(t, data, ex, 1)
	}
	t.Run("nil pointer", func(t *testing.T) {
		f(t, "null", ErrNilPointerReceiver, 0)
	})
	t.Run("eof", func(t *testing.T) {
		f2(t, "", io.EOF)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		f2(t, `true`, InvalidDigitError{})
	})
	t.Run("invalid null", func(t *testing.T) {
		f2(t, "nul", io.EOF)
	})
	t.Run("null", func(t *testing.T) {
		f2(t, "null", nil)
	})
	t.Run("valid", func(t *testing.T) {
		f2(t, "127", nil)
	})
}

func TestValDecoder_Native_Uint64(t *testing.T) {
	f := func(t *testing.T, data string, ex error, initValue uint64) {
		var p1 *uint64
		var p2 *uint64
		if initValue != 0 {
			b1 := initValue
			p1 = &b1
			b2 := initValue
			p2 = &b2
		}
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
	}
	f2 := func(t *testing.T, data string, ex error) {
		f(t, data, ex, 1)
	}
	t.Run("nil pointer", func(t *testing.T) {
		f(t, "null", ErrNilPointerReceiver, 0)
	})
	t.Run("eof", func(t *testing.T) {
		f2(t, "", io.EOF)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		f2(t, `true`, InvalidDigitError{})
	})
	t.Run("invalid null", func(t *testing.T) {
		f2(t, "nul", io.EOF)
	})
	t.Run("null", func(t *testing.T) {
		f2(t, "null", nil)
	})
	t.Run("valid", func(t *testing.T) {
		f2(t, "127", nil)
	})
}
