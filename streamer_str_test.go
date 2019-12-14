package jzon

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreamer_String(t *testing.T) {
	f := func(t *testing.T, s string) {
		v, err := json.Marshal(s)
		require.NoError(t, err)

		testStreamer(t, string(v), func(streamer *Streamer) {
			streamer.String(s)
		})
	}
	t.Run("empty", func(t *testing.T) {
		f(t, ``)
	})
	t.Run("quote", func(t *testing.T) {
		f(t, `"`)
	})
	t.Run("backslash", func(t *testing.T) {
		f(t, `\`)
	})
	t.Run("normal", func(t *testing.T) {
		f(t, "test")
	})
	t.Run("Line Feed", func(t *testing.T) {
		f(t, "\n")
	})
	t.Run("Carriage Return", func(t *testing.T) {
		f(t, "\r")
	})
	t.Run("tab", func(t *testing.T) {
		f(t, "\t")
	})
	t.Run("eof", func(t *testing.T) {
		f(t, "\u0000")
	})
	t.Run("utf8", func(t *testing.T) {
		f(t, "中")
	})
}
