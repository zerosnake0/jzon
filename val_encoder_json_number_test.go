package jzon

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValEncoder_JsonNumber(t *testing.T) {
	f := func(t *testing.T, n json.Number) {
		b, err := json.Marshal(n)
		require.NoError(t, err)
		testStreamer(t, string(b), func(s *Streamer) {
			s.Value(n)
		})
	}
	t.Run("empty", func(t *testing.T) {
		f(t, "")
	})
	t.Run("non empty", func(t *testing.T) {
		f(t, "-1.2e-3")
	})
	t.Run("invalid", func(t *testing.T) {
		// TODO:
	})
}
