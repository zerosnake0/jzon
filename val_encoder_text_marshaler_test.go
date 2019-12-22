package jzon

import (
	"encoding/json"
	"testing"

	"encoding"
	"github.com/stretchr/testify/require"
)

type testTextMarshaler struct {
	data []byte
	err  error
}

func (m testTextMarshaler) MarshalText() ([]byte, error) {
	return m.data, m.err
}

type testTextMarshaler2 struct {
	data []byte
	err  error
}

func (m *testTextMarshaler2) MarshalText() ([]byte, error) {
	return m.data, m.err
}

func TestValEncoder_TextMarshaler(t *testing.T) {
	f := func(t *testing.T, m encoding.TextMarshaler) {
		b, err := json.Marshal(m)
		require.NoError(t, err)
		testStreamer(t, string(b), func(s *Streamer) {
			s.Value(m)
		})
	}
	t.Run("non pointer receiver", func(t *testing.T) {
		t.Run("non pointer", func(t *testing.T) {
			f(t, testTextMarshaler{
				data: []byte(`{"a":1}`),
			})
		})
		t.Run("pointer", func(t *testing.T) {
			f(t, &testTextMarshaler{
				data: []byte(`{"a":2}`),
			})
		})
		t.Run("nil pointer", func(t *testing.T) {
			f(t, (*testTextMarshaler)(nil))
		})
	})
	t.Run("pointer receiver", func(t *testing.T) {
		t.Run("pointer", func(t *testing.T) {
			f(t, &testTextMarshaler2{
				data: []byte(`{"b":1}`),
			})
		})
		t.Run("nil pointer", func(t *testing.T) {
			f(t, (*testTextMarshaler2)(nil))
		})
	})
}
