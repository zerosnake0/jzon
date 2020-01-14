package jzon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func runeToAscii(s string) string {
	var ret string
	for _, r := range []rune(s) {
		if r < 128 {
			ret += string(r)
		} else {
			ret += fmt.Sprintf("\\u%04x", r)
		}
	}
	return ret
}

func testIteratorStr(t *testing.T, input string, ex error) {
	var s string
	err := json.Unmarshal([]byte(input), &s)
	require.Equalf(t, ex == nil, err == nil, "json.Marshal\nexp:%v\ngot:%v",
		ex, err)
	withIterator(input, func(it *Iterator) {
		s2, err := it.ReadString()
		checkError(t, ex, err)
		if ex == nil {
			require.Equalf(t, s, s2, "exp: %s\ngot: %s", s, s2)
			t.Logf("json: %s", s)
			t.Logf("jzon: %s", s2)
		}
	})
}

func TestIterator_Str_readU4(t *testing.T) {
	origin := `中`
	src := runeToAscii(`"` + origin + `"`)
	t.Run("eof", func(t *testing.T) {
		testIteratorStr(t, src[:3], io.EOF)
	})
	t.Run("eof2", func(t *testing.T) {
		testIteratorStr(t, src[:len(src)-1], io.EOF)
	})
	t.Run("invalid_unicode", func(t *testing.T) {
		testIteratorStr(t, `"\uG`, InvalidUnicodeCharError{})
	})
	t.Run("bytes", func(t *testing.T) {
		t.Run("", func(t *testing.T) {
			testIteratorStr(t, src, nil)
		})
		t.Run("", func(t *testing.T) {
			withIterator(src, func(it *Iterator) {
				s, err := it.ReadString()
				require.NoError(t, err)
				require.Equal(t, origin, s)
			})
		})
	})
	t.Run("reader", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			it.Reset(&oneByteReader{
				b: src,
			})
			s, err := it.ReadString()
			require.NoError(t, err)
			require.Equal(t, origin, s)
		})
	})
}

func TestIterator_Str_readEscapedChar(t *testing.T) {
	t.Run("eof", func(t *testing.T) {
		testIteratorStr(t, `"\`, io.EOF)
	})
	t.Run("control", func(t *testing.T) {
		data := `"\n"`
		t.Run("", func(t *testing.T) {
			testIteratorStr(t, data, nil)
		})
		t.Run("", func(t *testing.T) {
			withIterator(data, func(it *Iterator) {
				s, err := it.ReadString()
				require.NoError(t, err)
				require.Len(t, s, 1)
				require.Equal(t, "\n", s)
			})
		})
	})
	t.Run("invalid escape char", func(t *testing.T) {
		testIteratorStr(t, `"\a"`, InvalidEscapeCharError{})
	})
	t.Run("unicode error 1", func(t *testing.T) {
		testIteratorStr(t, `"\u0`, io.EOF)
	})
	t.Run("surrogate err 1", func(t *testing.T) {
		testIteratorStr(t, `"\ud800`, io.EOF)
	})
	t.Run("surrogate incomplete", func(t *testing.T) {
		data := `"\ud800"`
		t.Run("", func(t *testing.T) {
			testIteratorStr(t, data, nil)
		})
		t.Run("", func(t *testing.T) {
			withIterator(data, func(it *Iterator) {
				s, err := it.ReadString()
				require.NoError(t, err)
				require.Equal(t, string(runeError), s)
			})
		})
	})
	t.Run("surrogate err 2", func(t *testing.T) {
		testIteratorStr(t, `"\ud800\`, io.EOF)
	})
	t.Run("surrogate err 3", func(t *testing.T) {
		testIteratorStr(t, `"\ud800\a`, InvalidEscapeCharError{})
	})
	t.Run("surrogate other escaped char", func(t *testing.T) {
		data := `"\ud800\n"`
		t.Run("", func(t *testing.T) {
			testIteratorStr(t, data, nil)
		})
		t.Run("", func(t *testing.T) {
			withIterator(data, func(it *Iterator) {
				s, err := it.ReadString()
				require.NoError(t, err)
				require.Equal(t, string(runeError)+"\n", s)
			})
		})
	})
	t.Run("surrogate err 4", func(t *testing.T) {
		testIteratorStr(t, `"\ud800\u`, io.EOF)
	})
	t.Run("surrogate runeError", func(t *testing.T) {
		data := `"\udc00\u0000"`
		t.Run("", func(t *testing.T) {
			testIteratorStr(t, data, nil)
		})
		t.Run("", func(t *testing.T) {
			withIterator(data, func(it *Iterator) {
				s, err := it.ReadString()
				require.NoError(t, err)
				require.Equal(t, `"\ufffd\x00"`, fmt.Sprintf("%+q", s))
			})
		})
	})
	t.Run("surrogate", func(t *testing.T) {
		data := `"\uD852\uDF62"`
		t.Run("", func(t *testing.T) {
			testIteratorStr(t, data, nil)
		})
		t.Run("", func(t *testing.T) {
			withIterator(data, func(it *Iterator) {
				s, err := it.ReadString()
				require.NoError(t, err)
				require.Equal(t, "𤭢", s)
			})
		})
	})
	t.Run("non surrogate", func(t *testing.T) {
		s := `"\u4e2d"` // 中
		t.Run("", func(t *testing.T) {
			testIteratorStr(t, s, nil)
		})
		t.Run("", func(t *testing.T) {
			withIterator(s, func(it *Iterator) {
				s, err := it.ReadString()
				require.NoError(t, err)
				require.Equal(t, "中", s)
			})
		})
	})
	t.Run("surrogate consecutive", func(t *testing.T) {
		data := `"\udc00\uD852\uDF62"`
		t.Run("", func(t *testing.T) {
			testIteratorStr(t, data, nil)
		})
		withIterator(data, func(it *Iterator) {
			s, err := it.ReadString()
			require.NoError(t, err)
			require.Equal(t, "\ufffd𤭢", s)
			require.Equal(t, `"\ufffd\U00024b62"`, fmt.Sprintf("%+q", s))
		})
	})
}

func TestIterator_Str_readStringAsSlice(t *testing.T) {
	t.Run("reader error", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			e := errors.New("test")
			it.Reset(&oneByteReader{
				err: e,
			})
			_, err := it.ReadString()
			require.Error(t, e, err)
		})
	})
	t.Run("bad value type", func(t *testing.T) {
		withIterator(`1`, func(it *Iterator) {
			_, err := it.ReadString()
			require.IsType(t, UnexpectedByteError{}, err)
		})
	})
	t.Run("simple", func(t *testing.T) {
		withIterator(`"abc123"`, func(it *Iterator) {
			s, err := it.ReadString()
			require.NoError(t, err)
			require.Equal(t, "abc123", s)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("escape error", func(t *testing.T) {
		withIterator(`"\`, func(it *Iterator) {
			_, err := it.ReadString()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("escape", func(t *testing.T) {
		withIterator(`"\n"`, func(it *Iterator) {
			s, err := it.ReadString()
			require.NoError(t, err)
			require.Equal(t, "\n", s)
		})
	})
	t.Run("invalid string char", func(t *testing.T) {
		withIterator("\"\x00\"", func(it *Iterator) {
			_, err := it.ReadString()
			require.IsType(t, InvalidStringCharError{}, err)
		})
	})
	t.Run("reade err", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			it.Reset(&oneByteReader{
				b: `"\n`,
			})
			_, err := it.ReadString()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("reader", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			it.Reset(&oneByteReader{
				b: `"\n"`,
			})
			s, err := it.ReadString()
			require.NoError(t, err)
			require.Len(t, s, 1)
			require.Equal(t, "\n", s)
		})
	})
}

func TestIterator_Str_ReadStringAsSlice(t *testing.T) {
	t.Run("not string", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			_, err := it.ReadStringAndAppend(nil)
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("normal", func(t *testing.T) {
		withIterator(`"abc"`, func(it *Iterator) {
			buf := make([]byte, 0, 32)
			ret, err := it.ReadStringAndAppend(buf)
			require.NoError(t, err)
			p1 := (*sliceHeader)(unsafe.Pointer(&buf))
			p2 := (*sliceHeader)(unsafe.Pointer(&ret))
			require.Equal(t, p1.Data, p2.Data)
			require.Equal(t, p1.Cap, p2.Cap)
			require.Equal(t, []byte("abc"), ret)
		})
	})
}

func TestIterator_Str_appendRune(t *testing.T) {
	var b []byte

	b = appendRune(b[:0], 0)
	require.Equal(t, []byte{0}, b)

	b = appendRune(b[:0], 1<<7)
	require.Equal(t, []byte{0xc2, 0x80}, b)

	b = appendRune(b[:0], maxRune+1)
	require.Equal(t, []byte{0xef, 0xbf, 0xbd}, b)

	b = appendRune(b[:0], surrogateMin)
	require.Equal(t, []byte{0xef, 0xbf, 0xbd}, b)

	b = appendRune(b[:0], surrogateMax)
	require.Equal(t, []byte{0xef, 0xbf, 0xbd}, b)

	b = appendRune(b[:0], rune3Max)
	require.Equal(t, []byte{0xef, 0xbf, 0xbf}, b)

	b = appendRune(b[:0], rune3Max+1)
	require.Equal(t, []byte{0xf0, 0x90, 0x80, 0x80}, b)
}
