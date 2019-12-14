package jzon

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func testStreamerString(t *testing.T, s string) {
	v, err := json.Marshal(s)
	require.NoError(t, err)

	testStreamer(t, string(v), func(streamer *Streamer) {
		streamer.String(s)
	})
}

func TestStreamer_String(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		testStreamerString(t, ``)
	})
	t.Run("quote", func(t *testing.T) {
		testStreamerString(t, `"`)
	})
	t.Run("backslash", func(t *testing.T) {
		testStreamerString(t, `\`)
	})
	t.Run("normal", func(t *testing.T) {
		testStreamerString(t, "test")
	})
	t.Run("Line Feed", func(t *testing.T) {
		testStreamerString(t, "\n")
	})
	t.Run("Carriage Return", func(t *testing.T) {
		testStreamerString(t, "\r")
	})
	t.Run("tab", func(t *testing.T) {
		testStreamerString(t, "\t")
	})
	t.Run("eof", func(t *testing.T) {
		testStreamerString(t, "\u0000")
	})
	t.Run("utf8", func(t *testing.T) {
		testStreamerString(t, "中")
	})
}

func TestStreamer_String_Html(t *testing.T) {
	t.Run("bracket", func(t *testing.T) {
		testStreamerString(t, "<>")
	})
	t.Run("and", func(t *testing.T) {
		testStreamerString(t, "&")
	})
}
