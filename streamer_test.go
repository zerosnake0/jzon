package jzon

import (
	"bytes"
	"encoding/json"
	"errors"
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

func testStreamerWithEncoder(t *testing.T, enc *Encoder, exp string, cb func(s *Streamer)) {
	streamer := enc.NewStreamer()
	defer enc.ReturnStreamer(streamer)
	var b bytes.Buffer
	streamer.Reset(&b)

	cb(streamer)
	err := streamer.Flush()
	require.NoError(t, err)

	s := b.String()
	require.Equalf(t, exp, s, "expect %q but got %q", exp, s)
	t.Logf("got %q", s)
}

func testStreamer(t *testing.T, exp string, cb func(s *Streamer)) {
	testStreamerWithEncoder(t, DefaultEncoder, exp, cb)
}

func checkEncodeWithStandard(t *testing.T, enc *Encoder, obj interface{}, cb func(s *Streamer)) {
	buf, err := json.Marshal(obj)

	streamer := enc.NewStreamer()
	defer enc.ReturnStreamer(streamer)
	cb(streamer)

	if err != nil {
		require.Error(t, streamer.Error)
	} else {
		t.Logf("got %s", buf)
		require.NoError(t, streamer.Error)
		require.Equal(t, buf, streamer.buffer, "expecting %s but got %s",
			buf, streamer.buffer)
	}
}

func checkEncodeValueWithStandard(t *testing.T, enc *Encoder, obj interface{}) {
	checkEncodeWithStandard(t, enc, obj, func(s *Streamer) {
		s.Value(obj)
	})
}

func testStreamerChainError(t *testing.T, cb func(s *Streamer)) {
	s := DefaultEncoder.NewStreamer()
	defer DefaultEncoder.ReturnStreamer(s)

	var b bytes.Buffer
	s.Reset(&b)

	e := errors.New("test")
	s.Error = e
	cb(s)

	require.Equal(t, e, s.Error)
	require.Equal(t, e, s.Flush())
	require.Len(t, s.buffer, 0)
	require.Equal(t, 0, b.Len())
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
	t.Run("raw string", func(t *testing.T) {
		testStreamer(t, "abc", func(s *Streamer) {
			s.RawString("abc")
		})
	})
	t.Run("raw", func(t *testing.T) {
		testStreamer(t, "abc", func(s *Streamer) {
			s.Raw([]byte("abc"))
		})
	})
	t.Run("null", func(t *testing.T) {
		testStreamer(t, "null", func(s *Streamer) {
			s.Null()
		})
	})
	t.Run("true", func(t *testing.T) {
		testStreamer(t, "true", func(s *Streamer) {
			s.True()
		})
		testStreamer(t, "true", func(s *Streamer) {
			s.Bool(true)
		})
	})
	t.Run("false", func(t *testing.T) {
		testStreamer(t, "false", func(s *Streamer) {
			s.False()
		})
		testStreamer(t, "false", func(s *Streamer) {
			s.Bool(false)
		})
	})
	t.Run("array (empty)", func(t *testing.T) {
		testStreamer(t, "[]", func(s *Streamer) {
			s.ArrayStart().ArrayEnd()
		})
	})
	t.Run("array (nested 1)", func(t *testing.T) {
		count := 10
		s := strings.ReplaceAll(nestedArray1(count), " ", "")
		testStreamer(t, s, func(s *Streamer) {
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
		testStreamer(t, s, func(s *Streamer) {
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
		testStreamer(t, s, func(s *Streamer) {
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
		testStreamer(t, "{}", func(s *Streamer) {
			s.ObjectStart().ObjectEnd()
		})
	})
	t.Run("object (nested)", func(t *testing.T) {
		count := 5
		s := strings.ReplaceAll(nestedObject(count), " ", "")
		testStreamer(t, s, func(s *Streamer) {
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
		testStreamer(t, s, func(s *Streamer) {
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

func TestStreamer_ChainError(t *testing.T) {
	t.Run("raw string", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			s.RawString(`"test"`)
		})
	})
	t.Run("raw", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			s.Raw([]byte(`"test"`))
		})
	})
	t.Run("null", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			s.Null()
		})
	})
	t.Run("true", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			s.True()
		})
	})
	t.Run("false", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			s.False()
		})
	})
	t.Run("object start", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			s.ObjectStart()
		})
	})
	t.Run("object end", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			s.ObjectEnd()
		})
	})
	t.Run("field", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			s.Field("test")
		})
	})
	t.Run("array start", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			s.ArrayStart()
		})
	})
	t.Run("array end", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			s.ArrayEnd()
		})
	})
}