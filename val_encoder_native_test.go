package jzon

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestValEncoder_Bool(t *testing.T) {
	f := func(t *testing.T, b bool) {
		checkEncodeValueWithStandard(t, DefaultEncoder, b)
	}
	t.Run("true", func(t *testing.T) {
		f(t, true)
	})
	t.Run("false", func(t *testing.T) {
		f(t, false)
	})
	t.Run("nil ptr", func(t *testing.T) {
		checkEncodeValueWithStandard(t, DefaultEncoder, (*bool)(nil))
	})
	t.Run("ptr", func(t *testing.T) {
		b := true
		checkEncodeValueWithStandard(t, DefaultEncoder, &b)
	})
}

func TestValEncoder_Bool_Kind(t *testing.T) {
	type Bool bool
	f := func(t *testing.T, b Bool) {
		checkEncodeValueWithStandard(t, DefaultEncoder, b)
	}
	t.Run("true", func(t *testing.T) {
		f(t, true)
	})
	t.Run("false", func(t *testing.T) {
		f(t, false)
	})
	t.Run("nil ptr", func(t *testing.T) {
		checkEncodeValueWithStandard(t, DefaultEncoder, (*Bool)(nil))
	})
	t.Run("ptr", func(t *testing.T) {
		b := Bool(true)
		checkEncodeValueWithStandard(t, DefaultEncoder, &b)
	})
}

type testBoolEncoder struct{}

func (*testBoolEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if ptr == nil {
		s.Null()
		return
	}
	s.Bool(!*(*bool)(ptr))
}

func TestValEncoder_Bool_Kind_CustomEncoder(t *testing.T) {
	enc := NewEncoder(&EncoderOption{
		ValEncoders: map[reflect.Type]ValEncoder{
			reflect.TypeOf(true): (*testBoolEncoder)(nil),
		},
	})

	testStreamerWithEncoder(t, enc, "false", func(s *Streamer) {
		s.Value(true)
	})
	testStreamerWithEncoder(t, enc, "true", func(s *Streamer) {
		s.Value(false)
	})

	type Bool bool

	testStreamerWithEncoder(t, enc, "false", func(s *Streamer) {
		s.Value(Bool(true))
	})
	testStreamerWithEncoder(t, enc, "true", func(s *Streamer) {
		s.Value(Bool(false))
	})
}

func TestValEncoder_String(t *testing.T) {
	f := func(t *testing.T, str string) {
		checkEncodeValueWithStandard(t, DefaultEncoder, str)
	}
	t.Run("test", func(t *testing.T) {
		f(t, "test")
	})
	t.Run("escape", func(t *testing.T) {
		f(t, "\"\n\r\t\\")
	})
	t.Run("html escape", func(t *testing.T) {
		f(t, "<>&")
	})
	t.Run("nil ptr", func(t *testing.T) {
		checkEncodeValueWithStandard(t, DefaultEncoder, (*string)(nil))
	})
	t.Run("ptr", func(t *testing.T) {
		s := "test"
		checkEncodeValueWithStandard(t, DefaultEncoder, &s)
	})
}
