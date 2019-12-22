package jzon

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValEncoder_JsonRawMessage(t *testing.T) {
	f := func(t *testing.T, s string) {
		msg := json.RawMessage(s)
		b, err := json.Marshal(msg)
		require.NoError(t, err)
		testStreamer(t, string(b), func(s *Streamer) {
			s.Value(msg)
		})
	}
	t.Run("null", func(t *testing.T) {
		f(t, "null")
	})
}
