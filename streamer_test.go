package jzon

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type badWriter struct {
	data string
}

func (w *badWriter) Write(data []byte) (int, error) {
	n := len(data)
	if n > 0 {
		n -= 1
	}
	w.data = string(data)
	return n, nil
}

func TestStreamer_Flush(t *testing.T) {
	t.Run("no writer attached", func(t *testing.T) {
		streamer := NewStreamer()
		defer ReturnStreamer(streamer)
		err := streamer.Flush()
		require.Equal(t, NoWriterAttachedError, err)
	})
	t.Run("bad writer implementation", func(t *testing.T) {
		streamer := NewStreamer()
		defer ReturnStreamer(streamer)
		var (
			w   badWriter
			err error
		)
		streamer.Reset(&w)
		streamer.True()

		err = streamer.Flush()
		require.NoError(t, err)
		require.Equal(t, "true", w.data)

		err = streamer.Flush()
		require.NoError(t, err)
		require.Equal(t, "e", w.data)
	})
}

func TestStreamer(t *testing.T) {
	f := func(t *testing.T, exp string, cb func(s *Streamer)) {
		streamer := NewStreamer()
		defer ReturnStreamer(streamer)
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
	t.Run("object (nested)", func(t *testing.T) {
		count := 5
		s := strings.ReplaceAll(nestedObject(count), " ", "")
		f(t, s, func(s *Streamer) {
			for i := 0; i < count; i++ {
				s.ObjectStart().
					Field("a").ObjectStart().ObjectEnd().
					Field("b")
			}
			s.ObjectStart().ObjectEnd()
			for i := 0; i < count; i++ {
				s.ObjectEnd()
			}
		})
	})
	t.Run("object (nested with array)", func(t *testing.T) {
		count := 5
		s := strings.ReplaceAll(nestedObjectWithArray(count), " ", "")
		f(t, s, func(s *Streamer) {
			for i := 0; i < count; i++ {
				s.ObjectStart().
					Field("a").ArrayStart().ArrayEnd().
					Field("b")
			}
			s.ArrayStart().ArrayEnd()
			for i := 0; i < count; i++ {
				s.ObjectEnd()
			}
		})
	})
}
