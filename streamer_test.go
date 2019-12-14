package jzon

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreamer(t *testing.T) {
	f := func(t *testing.T, exp string, cb func(s *Streamer)) {
		streamer := NewStreamer()
		b := bytes.NewBuffer(nil)
		streamer.Reset(b)
		cb(streamer)
		err := streamer.Flush()
		require.NoError(t, err)
		s := b.String()
		require.Equalf(t, exp, s, "expect %q but got %q", exp, s)
		t.Logf("got %q", s)
	}
	t.Run("raw string", func(t *testing.T) {
		f(t, "abc", func(s *Streamer) {
			s.RawString("abc")
		})
	})
	t.Run("raw", func(t *testing.T) {
		f(t, "abc", func(s *Streamer) {
			s.Raw([]byte("abc"))
		})
	})
	t.Run("null", func(t *testing.T) {
		f(t, "null", func(s *Streamer) {
			s.Null()
		})
	})
	t.Run("true", func(t *testing.T) {
		f(t, "true", func(s *Streamer) {
			s.True()
		})
		f(t, "true", func(s *Streamer) {
			s.Bool(true)
		})
	})
	t.Run("false", func(t *testing.T) {
		f(t, "false", func(s *Streamer) {
			s.False()
		})
		f(t, "false", func(s *Streamer) {
			s.Bool(false)
		})
	})
	t.Run("array (empty)", func(t *testing.T) {
		f(t, "[]", func(s *Streamer) {
			s.ArrayStart().ArrayEnd()
		})
	})
	t.Run("array (nested 1)", func(t *testing.T) {
		count := 10
		s := strings.ReplaceAll(nestedArray1(count), " ", "")
		f(t, s, func(s *Streamer) {
			for i := 0; i < count; i++ {
				s.ArrayStart()
			}
			s.ArrayStart().ArrayEnd()
			for i := 0; i < count; i++ {
				s.ArrayEnd()
			}
		})
	})
	t.Run("array (nested 2)", func(t *testing.T) {
		count := 10
		s := strings.ReplaceAll(nestedArray2(count), " ", "")
		f(t, s, func(s *Streamer) {
			for i := 0; i < count; i++ {
				s.ArrayStart().
					ArrayStart().ArrayEnd()
			}
			s.ArrayStart().ArrayEnd()
			for i := 0; i < count; i++ {
				s.ArrayEnd()
			}
		})
	})
	t.Run("array (nested with object)", func(t *testing.T) {
		count := 10
		s := strings.ReplaceAll(nestedArrayWithObject(count), " ", "")
		f(t, s, func(s *Streamer) {
			for i := 0; i < count; i++ {
				s.ArrayStart().
					ObjectStart().ObjectEnd()
			}
			s.ArrayStart().ArrayEnd()
			for i := 0; i < count; i++ {
				s.ArrayEnd()
			}
		})
	})
	t.Run("object (empty)", func(t *testing.T) {
		f(t, "{}", func(s *Streamer) {
			s.ObjectStart().ObjectEnd()
		})
	})
	// t.Run("object (1 item)", func(t *testing.T) {
	// 	f(t, "{}", func(s *Streamer) {
	// 		s.ObjectStart().ObjectEnd()
	// 	})
	// })
}
