package jzon

import (
	"io"
	"testing"
)

func TestValDecoder_Native_Int(t *testing.T) {
	f := func(t *testing.T, data string, ex error, initValue int) {
		var p1 *int
		var p2 *int
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

func TestValDecoder_Native_Int8(t *testing.T) {
	f := func(t *testing.T, data string, ex error, initValue int8) {
		var p1 *int8
		var p2 *int8
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

func TestValDecoder_Native_Int16(t *testing.T) {
	f := func(t *testing.T, data string, ex error, initValue int16) {
		var p1 *int16
		var p2 *int16
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

func TestValDecoder_Native_Int32(t *testing.T) {
	f := func(t *testing.T, data string, ex error, initValue int32) {
		var p1 *int32
		var p2 *int32
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

func TestValDecoder_Native_Int64(t *testing.T) {
	f := func(t *testing.T, data string, ex error, initValue int64) {
		var p1 *int64
		var p2 *int64
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
