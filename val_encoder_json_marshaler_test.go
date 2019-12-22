package jzon

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type testJsonMarshaler struct {
	data []byte
	err  error
}

func (m testJsonMarshaler) MarshalJSON() ([]byte, error) {
	return m.data, m.err
}

type testJsonMarshaler2 struct {
	data []byte
	err  error
}

func (m *testJsonMarshaler2) MarshalJSON() ([]byte, error) {
	return m.data, m.err
}

func TestValEncoder_JsonMarshaler(t *testing.T) {
	f := func(t *testing.T, m json.Marshaler) {
		b, err := json.Marshal(m)
		require.NoError(t, err)
		testStreamer(t, string(b), func(s *Streamer) {
			s.Value(m)
		})
	}
	t.Run("non pointer", func(t *testing.T) {
		f(t, testJsonMarshaler{
			data: []byte(`{"a":1}`),
		})
	})
	t.Run("pointer", func(t *testing.T) {
		f(t, &testJsonMarshaler2{
			data: []byte(`{"b":2}`),
		})
	})
}
