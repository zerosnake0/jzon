package jzon

import (
	"encoding/json"
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func TestValEncoder_Bool(t *testing.T) {
	f := func(t *testing.T, b bool) {
		checkEncodeValueWithStandard(t, b, nil)
	}
	t.Run("true", func(t *testing.T) {
		f(t, true)
	})
	t.Run("false", func(t *testing.T) {
		f(t, false)
	})
	t.Run("nil ptr", func(t *testing.T) {
		checkEncodeValueWithStandard(t, (*bool)(nil), nil)
	})
	t.Run("ptr", func(t *testing.T) {
		b := true
		checkEncodeValueWithStandard(t, &b, nil)
	})
}

func TestValEncoder_Bool_Kind(t *testing.T) {
	type Bool bool
	f := func(t *testing.T, b Bool) {
		checkEncodeValueWithStandard(t, b, nil)
	}
	t.Run("true", func(t *testing.T) {
		f(t, true)
	})
	t.Run("false", func(t *testing.T) {
		f(t, false)
	})
	t.Run("nil ptr", func(t *testing.T) {
		checkEncodeValueWithStandard(t, (*Bool)(nil), nil)
	})
	t.Run("ptr", func(t *testing.T) {
		b := Bool(true)
		checkEncodeValueWithStandard(t, &b, nil)
	})
}

type testBoolEncoder struct{}

func (*testBoolEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
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
		checkEncodeValueWithStandard(t, str, nil)
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
		checkEncodeValueWithStandard(t, (*string)(nil), nil)
	})
	t.Run("ptr", func(t *testing.T) {
		s := "test"
		checkEncodeValueWithStandard(t, &s, nil)
	})
	t.Run("unicode", func(t *testing.T) {
		s := "\xe6\x97\xa5\xe6\x9c\xac\xff\xaa\x9e"
		checkEncodeValueWithStandard(t, s, nil)
	})
	t.Run("invalid unicode", func(t *testing.T) {
		s := `"invalid: \uD834x\uDD1E"`
		var s2 string
		err := json.Unmarshal([]byte(s), &s2)
		require.NoError(t, err)
		b, err := json.Marshal(s2)
		require.NoError(t, err)
		checkEncodeValueWithStandard(t, string(b), nil)
	})
}
