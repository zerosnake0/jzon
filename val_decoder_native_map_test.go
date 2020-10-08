package jzon

import (
	"io"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValDecoder_Native_Map(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
	}
	t.Run("nil pointer", func(t *testing.T) {
		f(t, "null", NilPointerReceiverError, nil, nil)
	})
	t.Run("eof", func(t *testing.T) {
		m1 := map[string]int{"a": 1}
		m2 := map[string]int{"a": 1}
		f(t, "", io.EOF, &m1, &m2)
	})
	t.Run("invalid null", func(t *testing.T) {
		m1 := map[string]int{"a": 1}
		m2 := map[string]int{"a": 1}
		f(t, "nul", io.EOF, &m1, &m2)
	})
	t.Run("null", func(t *testing.T) {
		m1 := map[string]int{"a": 1}
		m2 := map[string]int{"a": 1}
		f(t, "null", nil, &m1, &m2)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		m1 := map[string]int{"a": 1}
		m2 := map[string]int{"a": 1}
		f(t, " } ", UnexpectedByteError{}, &m1, &m2)
	})
	t.Run("eof after bracket", func(t *testing.T) {
		m1 := map[string]int{"a": 1}
		m2 := map[string]int{"a": 1}
		f(t, " { ", io.EOF, &m1, &m2)
	})
	t.Run("empty", func(t *testing.T) {
		m1 := map[string]int{"a": 1}
		m2 := map[string]int{"a": 1}
		f(t, " { } ", nil, &m1, &m2)
	})
	t.Run("empty on nil", func(t *testing.T) {
		var m1 map[string]int
		var m2 map[string]int
		f(t, " { } ", nil, &m1, &m2)
	})
	t.Run("value on nil", func(t *testing.T) {
		var m1 map[string]int
		var m2 map[string]int
		f(t, ` { "a" : 1 } `, nil, &m1, &m2)
	})
	t.Run("bad key", func(t *testing.T) {
		m1 := map[string]int{"a": 1}
		m2 := map[string]int{"a": 1}
		f(t, ` { "a`, io.EOF, &m1, &m2)
	})
	t.Run("eof after key", func(t *testing.T) {
		m1 := map[string]int{"a": 1}
		m2 := map[string]int{"a": 1}
		f(t, ` { "a" `, io.EOF, &m1, &m2)
	})
	t.Run("invalid colon", func(t *testing.T) {
		m1 := map[string]int{"a": 1}
		m2 := map[string]int{"a": 1}
		f(t, ` { "a" } `, UnexpectedByteError{}, &m1, &m2)
	})
	t.Run("bad value", func(t *testing.T) {
		m1 := map[string]int{"a": 1}
		m2 := map[string]int{"a": 1}
		f(t, ` { "b" : "c" } `, InvalidDigitError{}, &m1, &m2)
	})
	t.Run("eof after value", func(t *testing.T) {
		m1 := map[string]int{"a": 1}
		m2 := map[string]int{"a": 1}
		f(t, ` { "b" : 2 `, io.EOF, &m1, &m2)
	})
	t.Run("bad comma", func(t *testing.T) {
		m1 := map[string]int{"a": 1}
		m2 := map[string]int{"a": 1}
		f(t, ` { "b" : 2 { `, UnexpectedByteError{}, &m1, &m2)
	})
	t.Run("more items", func(t *testing.T) {
		m1 := map[string]int{"a": 1}
		m2 := map[string]int{"a": 1}
		f(t, ` { "b" : 2 , "c" : 3 } `, nil, &m1, &m2)
	})
	debug.FreeOSMemory()
}

type testMapIntKey int

func (k *testMapIntKey) UnmarshalText(data []byte) error {
	*k = testMapIntKey(len(data))
	return nil
}

type testMapStringKey string

func (k *testMapStringKey) UnmarshalText(data []byte) error {
	*k = testMapStringKey("`" + string(data) + "`")
	return nil
}

func TestValDecoder_Native_Map_KeyDecoder_TextUnmarshaler(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
	}
	t.Run("not supported", func(t *testing.T) {
		type key testTextUnmarshaler
		m1 := map[key]int{{
			data: "a",
		}: 1}
		m2 := map[key]int{{
			data: "a",
		}: 1}
		f(t, ` { "b" : 2 } `, TypeNotSupportedError(""), &m1, &m2)
	})
	t.Run("string", func(t *testing.T) {
		type key = testTextUnmarshaler
		m1 := map[key]int{{
			data: "a",
		}: 1}
		m2 := map[key]int{{
			data: "a",
		}: 1}
		f(t, ` { "b" : 2 } `, nil, &m1, &m2)
	})
	t.Run("int key", func(t *testing.T) {
		m1 := map[testMapIntKey]testMapIntKey{testMapIntKey(1): 2}
		m2 := map[testMapIntKey]testMapIntKey{testMapIntKey(1): 2}
		f(t, ` { "3" : "4" } `, nil, &m1, &m2)
	})
	t.Run("string key", func(t *testing.T) {
		// the UnmarshalText of the key is ignored
		m1 := map[testMapStringKey]testMapStringKey{testMapStringKey("1"): "2"}
		m2 := map[testMapStringKey]testMapStringKey{testMapStringKey("1"): "2"}
		v := "go1.13.15"
		if goVersion.LessEqual(v) {
			err := Unmarshal([]byte(` { "3" : "4" } `), &m2)
			require.NoError(t, err)
			require.Equal(t, testMapStringKey("`4`"), m2["`3`"])
		} else {
			f(t, ` { "3" : "4" } `, nil, &m1, &m2)
		}
	})
}

func TestValDecoder_Native_Map_KeyDecoder_String(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
	}
	t.Run("string", func(t *testing.T) {
		type key string
		m1 := map[key]int{"a": 1}
		m2 := map[key]int{"a": 1}
		f(t, ` { "b" : 2 } `, nil, &m1, &m2)
	})
	t.Run("string 2", func(t *testing.T) {
		type key string
		m1 := map[key]int{"a": 1}
		m2 := map[key]int{"a": 1}
		f(t, ` { "\"" : 2 } `, nil, &m1, &m2)
	})
}

func TestValDecoder_Native_Map_KeyDecoder_Int(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
	}
	t.Run("int8", func(t *testing.T) {
		type key int8
		t.Run("invalid", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { "b" : 3 } `, InvalidDigitError{}, &m1, &m2)
		})
		t.Run("no leading quote", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { 2`, UnexpectedByteError{}, &m1, &m2)
		})
		t.Run("no trimming quote", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { "2`, io.EOF, &m1, &m2)
		})
		t.Run("overflow", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { "128" : 3 } `, IntOverflowError{}, &m1, &m2)
		})
		t.Run("leading space", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { " 2" : 3 } `, InvalidDigitError{}, &m1, &m2)
		})
		t.Run("trailing space", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { "2 " : 3 } `, UnexpectedByteError{}, &m1, &m2)
		})
		t.Run("valid", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { "2" : 3 } `, nil, &m1, &m2)
		})
	})
	t.Run("int16", func(t *testing.T) {
		type key int16
		t.Run("invalid", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { "b" : 3 } `, InvalidDigitError{}, &m1, &m2)
		})
		t.Run("no leading quote", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { 2`, UnexpectedByteError{}, &m1, &m2)
		})
		t.Run("no trimming quote", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { "2`, io.EOF, &m1, &m2)
		})
		t.Run("overflow", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { "32768" : 3 } `, IntOverflowError{}, &m1, &m2)
		})
		t.Run("leading space", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { " 2" : 3 } `, InvalidDigitError{}, &m1, &m2)
		})
		t.Run("trailing space", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { "2 " : 3 } `, UnexpectedByteError{}, &m1, &m2)
		})
		t.Run("valid", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { "2" : 3 } `, nil, &m1, &m2)
		})
	})
	t.Run("int32", func(t *testing.T) {
		type key int32
		t.Run("invalid", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { "b" : 3 } `, InvalidDigitError{}, &m1, &m2)
		})
		t.Run("no leading quote", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { 2`, UnexpectedByteError{}, &m1, &m2)
		})
		t.Run("no trimming quote", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { "2`, io.EOF, &m1, &m2)
		})
		t.Run("overflow", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { "2147483649" : 3 } `, IntOverflowError{}, &m1, &m2)
		})
		t.Run("leading space", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { " 2" : 3 } `, InvalidDigitError{}, &m1, &m2)
		})
		t.Run("trailing space", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { "2 " : 3 } `, UnexpectedByteError{}, &m1, &m2)
		})
		t.Run("valid", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { "2" : 3 } `, nil, &m1, &m2)
		})
	})
	t.Run("int64", func(t *testing.T) {
		type key int64
		t.Run("invalid", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { "b" : 3 } `, InvalidDigitError{}, &m1, &m2)
		})
		t.Run("no leading quote", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { 2`, UnexpectedByteError{}, &m1, &m2)
		})
		t.Run("no trimming quote", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { "2`, io.EOF, &m1, &m2)
		})
		t.Run("overflow", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { "9223372036854775808" : 3 } `, IntOverflowError{}, &m1, &m2)
		})
		t.Run("leading space", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { " 2" : 3 } `, InvalidDigitError{}, &m1, &m2)
		})
		t.Run("trailing space", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { "2 " : 3 } `, UnexpectedByteError{}, &m1, &m2)
		})
		t.Run("valid", func(t *testing.T) {
			m1 := map[key]int{1: 2}
			m2 := map[key]int{1: 2}
			f(t, ` { "2" : 3 } `, nil, &m1, &m2)
		})
	})
}

func TestValDecoder_Native_Map_KeyDecoder_Uint(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
	}
	t.Run("uint8", func(t *testing.T) {
		type key uint8
		t.Run("invalid", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { "b" : 3 } `, InvalidDigitError{}, &m1, &m2)
		})
		t.Run("no leading quote", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { 2`, UnexpectedByteError{}, &m1, &m2)
		})
		t.Run("no trimming quote", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { "2`, io.EOF, &m1, &m2)
		})
		t.Run("overflow", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { "256" : 3 } `, IntOverflowError{}, &m1, &m2)
		})
		t.Run("leading space", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { " 2" : 3 } `, InvalidDigitError{}, &m1, &m2)
		})
		t.Run("trailing space", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { "2 " : 3 } `, UnexpectedByteError{}, &m1, &m2)
		})
		t.Run("valid", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { "2" : 3 } `, nil, &m1, &m2)
		})
	})
	t.Run("uint16", func(t *testing.T) {
		type key uint16
		t.Run("invalid", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { "b" : 3 } `, InvalidDigitError{}, &m1, &m2)
		})
		t.Run("no leading quote", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { 2`, UnexpectedByteError{}, &m1, &m2)
		})
		t.Run("no trimming quote", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { "2`, io.EOF, &m1, &m2)
		})
		t.Run("overflow", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { "65536" : 3 } `, IntOverflowError{}, &m1, &m2)
		})
		t.Run("leading space", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { " 2" : 3 } `, InvalidDigitError{}, &m1, &m2)
		})
		t.Run("trailing space", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { "2 " : 3 } `, UnexpectedByteError{}, &m1, &m2)
		})
		t.Run("valid", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { "2" : 3 } `, nil, &m1, &m2)
		})
	})
	t.Run("uint32", func(t *testing.T) {
		type key uint32
		t.Run("invalid", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { "b" : 3 } `, InvalidDigitError{}, &m1, &m2)
		})
		t.Run("no leading quote", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { 2`, UnexpectedByteError{}, &m1, &m2)
		})
		t.Run("no trimming quote", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { "2`, io.EOF, &m1, &m2)
		})
		t.Run("overflow", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { "4294967296" : 3 } `, IntOverflowError{}, &m1, &m2)
		})
		t.Run("leading space", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { " 2" : 3 } `, InvalidDigitError{}, &m1, &m2)
		})
		t.Run("trailing space", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { "2 " : 3 } `, UnexpectedByteError{}, &m1, &m2)
		})
		t.Run("valid", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { "2" : 3 } `, nil, &m1, &m2)
		})
	})
	t.Run("uint64", func(t *testing.T) {
		type key uint64
		t.Run("invalid", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { "b" : 3 } `, InvalidDigitError{}, &m1, &m2)
		})
		t.Run("no leading quote", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { 2`, UnexpectedByteError{}, &m1, &m2)
		})
		t.Run("no trimming quote", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { "2`, io.EOF, &m1, &m2)
		})
		t.Run("overflow", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { "18446744073709551616" : 3 } `, IntOverflowError{}, &m1, &m2)
		})
		t.Run("leading space", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { " 2" : 3 } `, InvalidDigitError{}, &m1, &m2)
		})
		t.Run("trailing space", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { "2 " : 3 } `, UnexpectedByteError{}, &m1, &m2)
		})
		t.Run("valid", func(t *testing.T) {
			m1 := map[key]uint{1: 2}
			m2 := map[key]uint{1: 2}
			f(t, ` { "2" : 3 } `, nil, &m1, &m2)
		})
	})
}
