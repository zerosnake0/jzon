package jzon

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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

func testStreamerWithEncoderConfig(t *testing.T, encCfg *EncoderConfig, exp string, cb func(s *Streamer)) {
	streamer := encCfg.NewStreamer()
	defer streamer.Release()
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
	testStreamerWithEncoderConfig(t, DefaultEncoderConfig, exp, cb)
}

func jsonMarshal(o interface{}, jsonOpt func(*json.Encoder)) (_ []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			if err == nil {
				err = fmt.Errorf("panic: %v", e)
			}
		}
	}()
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if jsonOpt != nil {
		jsonOpt(enc)
	}
	if err = enc.Encode(o); err != nil {
		return
	}
	return buf.Bytes(), nil
}

func jsonEqual(s1, s2 []byte) (bool, error) {
	var err error

	s1 = bytes.TrimSpace(s1)
	s2 = bytes.TrimSpace(s2)

	switch s1[0] {
	case 'n':
		return "null" == localByteToString(s1) &&
			"null" == localByteToString(s2), nil
	case 't':
		return "true" == localByteToString(s1) &&
			"true" == localByteToString(s2), nil
	case 'f':
		return "false" == localByteToString(s1) &&
			"false" == localByteToString(s2), nil
	case '[':
		var arr1 []json.RawMessage
		if err = json.Unmarshal(s1, &arr1); err != nil {
			return false, err
		}
		var arr2 []json.RawMessage
		if err = json.Unmarshal(s2, &arr2); err != nil {
			return false, err
		}
		l := len(arr1)
		if l != len(arr2) {
			return false, nil
		}
		for i := 0; i < l; i++ {
			b, err := jsonEqual(arr1[i], arr2[i])
			if err != nil || !b {
				return b, err
			}
		}
		return true, nil
	case '{':
		var m1 map[string]json.RawMessage
		if err = json.Unmarshal(s1, &m1); err != nil {
			return false, err
		}
		var m2 map[string]json.RawMessage
		if err = json.Unmarshal(s2, &m2); err != nil {
			return false, err
		}
		l := len(m1)
		if l != len(m2) {
			return false, nil
		}
		for k := range m1 {
			b, err := jsonEqual(m1[k], m2[k])
			if err != nil || !b {
				return b, err
			}
		}
		return true, nil
	case '"':
		var str1 string
		if err = json.Unmarshal(s1, &str1); err != nil {
			return false, err
		}
		var str2 string
		if err = json.Unmarshal(s2, &str2); err != nil {
			return false, err
		}
		return str1 == str2, nil
	default:
		var f1 float64
		if err = json.Unmarshal(s1, &f1); err != nil {
			return false, err
		}
		var f2 float64
		if err = json.Unmarshal(s2, &f2); err != nil {
			return false, err
		}
		return f1 == f2, nil
	}
}

var (
	testNoEscapeEncoderConfig = NewEncoderConfig(&EncoderOption{
		EscapeHTML: false,
	})
)

func checkEncodeWithStandard(t *testing.T, obj interface{}, cb func(s *Streamer),
	expErr interface{}) {
	checkEncodeWithStandardInternal(t, nil, DefaultEncoderConfig, obj, cb, expErr)
	checkEncodeWithStandardInternal(t, func(encoder *json.Encoder) {
		encoder.SetEscapeHTML(false)
	}, testNoEscapeEncoderConfig, obj, cb, expErr)
}

func checkEncodeWithStandardInternal(t *testing.T, jsonOpt func(*json.Encoder), encCfg *EncoderConfig, obj interface{},
	cb func(s *Streamer), expErr interface{}) {
	buf, err := jsonMarshal(obj, jsonOpt)
	require.Equal(t, expErr == nil, err == nil, "json.Marshal\nexp: %v\ngot: %v",
		expErr, err)

	streamer := encCfg.NewStreamer()
	defer streamer.Release()
	func() {
		defer func() {
			if e := recover(); e != nil {
				if streamer.Error != nil {
					panic(e)
				}
				ex, ok := e.(error)
				if !ok {
					panic(e)
				}
				streamer.Error = ex
			}
		}()
		cb(streamer)
	}()

	if err != nil {
		t.Logf("json err: %v", err)
		t.Logf("jzon err: %v", streamer.Error)
		switch x := expErr.(type) {
		case reflect.Type:
			gotErrType := reflect.TypeOf(streamer.Error)
			if x.Kind() == reflect.Interface {
				require.True(t, gotErrType.Implements(x), "exp err:%v\ngot err:%v",
					x, streamer.Error)
			} else {
				require.Equal(t, x, gotErrType, "exp err:%v\ngot err:%v",
					x, streamer.Error)
			}
		case error:
			checkError(t, x, streamer.Error)
			// if reflect.TypeOf(errors.New("")) == reflect.TypeOf(expErr) {
			// 	require.Equalf(t, expErr, streamer.Error, "exp err:%v\ngot err:%v",
			// 		expErr, streamer.Error)
			// } else {
			// 	require.IsTypef(t, expErr, streamer.Error, "exp err:%v\ngot err:%v",
			// 		expErr, streamer.Error)
			// }
		}
		require.Error(t, streamer.Error, "json.Marshal error: %v", err)
	} else {
		t.Logf("got %s", buf)
		require.NoError(t, streamer.Error)
		b, err := jsonEqual(buf, streamer.buffer)
		require.NoErrorf(t, err, "final result\njson %s\njzon %s",
			bytes.TrimSpace(buf), bytes.TrimSpace(streamer.buffer))
		require.Truef(t, b, "final result\njson %s\njzon %s",
			bytes.TrimSpace(buf), bytes.TrimSpace(streamer.buffer))
	}
}

func checkEncodeValueWithStandard(t *testing.T, obj interface{}, expErr interface{}) {
	checkEncodeWithStandard(t, obj, func(s *Streamer) {
		s.Value(obj)
	}, expErr)
}

func testStreamerChainError(t *testing.T, cb func(s *Streamer)) {
	s := DefaultEncoderConfig.NewStreamer()
	defer s.Release()

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
		defer streamer.Release()
		err := streamer.Flush()
		require.Equal(t, NoWriterAttachedError, err)
	})
	t.Run("bad writer implementation", func(t *testing.T) {
		streamer := NewStreamer()
		defer streamer.Release()
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
	t.Run("array", func(t *testing.T) {
		t.Run("empty", func(t *testing.T) {
			testStreamer(t, "[]", func(s *Streamer) {
				s.ArrayStart().ArrayEnd()
			})
		})
		t.Run("nested 1", func(t *testing.T) {
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
		t.Run("nested 2", func(t *testing.T) {
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
		t.Run("nested with object", func(t *testing.T) {
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
	})
	t.Run("object", func(t *testing.T) {
		t.Run("empty", func(t *testing.T) {
			testStreamer(t, "{}", func(s *Streamer) {
				s.ObjectStart().ObjectEnd()
			})
		})
		t.Run("nested", func(t *testing.T) {
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
		t.Run("nested with array", func(t *testing.T) {
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
