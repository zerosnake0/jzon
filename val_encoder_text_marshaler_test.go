package jzon

import (
	"encoding"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValEncoder_TextMarshaler_ChainError(t *testing.T) {
	t.Run("direct", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			textMarshalerEncoder(0).Encode(nil, s)
		})
	})
	t.Run("dynamic", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			(*dynamicTextMarshalerEncoder)(nil).Encode(nil, s)
		})
	})
}

type testTextMarshaler struct {
	data string
	err  error
}

func (m testTextMarshaler) MarshalText() ([]byte, error) {
	return []byte(m.data), m.err
}

type testTextMarshaler2 struct {
	data string
	err  error
}

func (m *testTextMarshaler2) MarshalText() ([]byte, error) {
	return []byte(m.data), m.err
}

func TestValEncoder_TextMarshaler(t *testing.T) {
	f := func(t *testing.T, m encoding.TextMarshaler, err error) {
		checkEncodeValueWithStandard(t, DefaultEncoder, m, err)
	}
	t.Run("non pointer receiver", func(t *testing.T) {
		t.Run("non pointer", func(t *testing.T) {
			f(t, testTextMarshaler{
				data: `{"a":1}`,
			}, nil)
		})
		t.Run("non pointer error", func(t *testing.T) {
			e := errors.New("test")
			f(t, testTextMarshaler{
				data: `{"a":1}`,
				err:  e,
			}, e)
		})
		t.Run("pointer", func(t *testing.T) {
			f(t, &testTextMarshaler{
				data: `{"a":2}`,
			}, nil)
		})
		t.Run("pointer error", func(t *testing.T) {
			e := errors.New("test")
			f(t, &testTextMarshaler{
				data: `{"a":2}`,
				err:  e,
			}, e)
		})
		t.Run("nil pointer", func(t *testing.T) {
			f(t, (*testTextMarshaler)(nil), nil)
		})
	})
	t.Run("pointer receiver", func(t *testing.T) {
		t.Run("pointer", func(t *testing.T) {
			f(t, &testTextMarshaler2{
				data: `{"b":1}`,
			}, nil)
		})
		t.Run("pointer error", func(t *testing.T) {
			e := errors.New("test")
			f(t, &testTextMarshaler2{
				data: `{"b":1}`,
				err:  e,
			}, e)
		})
		t.Run("nil pointer", func(t *testing.T) {
			f(t, (*testTextMarshaler2)(nil), nil)
		})
	})
}

func TestValEncoder_DynamicTextMarshaler(t *testing.T) {
	t.Run("marshaler (nil)", func(t *testing.T) {
		// TODO: This test should be automatically fixed in the future golang version
		v := "go1.13.5"
		if goVersion.LessEqual(v) {
			var i encoding.TextMarshaler
			b, err := DefaultEncoder.Marshal(&i)
			require.NoError(t, err)
			require.Equal(t, "null", string(b))
		} else {
			var i encoding.TextMarshaler
			checkEncodeValueWithStandard(t, DefaultEncoder, &i, nil)
		}
	})
	t.Run("marshaler error", func(t *testing.T) {
		e := errors.New("test")
		var i encoding.TextMarshaler = testTextMarshaler{
			data: `"test"`,
			err:  e,
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, &i, e)
	})
	t.Run("marshaler", func(t *testing.T) {
		var i encoding.TextMarshaler = testTextMarshaler{
			data: `"test"`,
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, &i, nil)
	})
	t.Run("marshaler 2", func(t *testing.T) {
		var i encoding.TextMarshaler = &testTextMarshaler{
			data: `"test 2"`,
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, &i, nil)
	})
}
