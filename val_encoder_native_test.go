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
		buf, err := json.Marshal(b)
		require.NoError(t, err)
		testStreamer(t, string(buf), func(s *Streamer) {
			s.Value(b)
		})
	}
	t.Run("true", func(t *testing.T) {
		f(t, true)
	})
	t.Run("false", func(t *testing.T) {
		f(t, false)
	})
}

func TestValEncoder_Bool_Kind(t *testing.T) {
	type Bool bool
	f := func(t *testing.T, b Bool) {
		buf, err := json.Marshal(b)
		require.NoError(t, err)
		testStreamer(t, string(buf), func(s *Streamer) {
			s.Value(b)
		})
	}
	t.Run("true", func(t *testing.T) {
		f(t, true)
	})
	t.Run("false", func(t *testing.T) {
		f(t, false)
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
		buf, err := json.Marshal(str)
		require.NoError(t, err)
		testStreamer(t, string(buf), func(s *Streamer) {
			s.Value(str)
		})
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
}
