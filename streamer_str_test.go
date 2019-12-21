package jzon

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

var (
	noEscapeEncoder = NewEncoder(&EncoderOption{
		EscapeHTML: false,
	})
)

func TestStreamer_String_ChainError(t *testing.T) {
	testStreamerChainError(t, func(s *Streamer) {
		s.String("test")
	})
}

func testStreamerStringEscape(t *testing.T, s string, escape bool) {
	var (
		enc *Encoder
		b   bytes.Buffer
	)
	jsEnc := json.NewEncoder(&b)
	jsEnc.SetEscapeHTML(escape)
	jsEnc.Encode(s)
	if escape {
		enc = DefaultEncoder
	} else {
		enc = noEscapeEncoder
	}
	// json.Encoder will add a newline at the end
	exp := strings.TrimSpace(b.String())
	testStreamerWithEncoder(t, enc, exp, func(streamer *Streamer) {
		streamer.String(s)
	})
}

func testStreamerString(t *testing.T, s string) {
	testStreamerStringEscape(t, s, true)
	testStreamerStringEscape(t, s, false)
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
	t.Run("rune error", func(t *testing.T) {
		testStreamerString(t, "\xFF")
	})
	t.Run("2028", func(t *testing.T) {
		testStreamerString(t, "\u2028")
	})
	t.Run("2029", func(t *testing.T) {
		testStreamerString(t, "\u2029")
	})
}
