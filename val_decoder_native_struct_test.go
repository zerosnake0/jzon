package jzon

import (
	"io"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValDecoder_Native_Struct_Zero_Field(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		checkDecodeWithStandard(t, DefaultDecoder, data, ex, p1, p2)
	}
	t.Run("nil receiver", func(t *testing.T) {
		f(t, "null", NilPointerReceiverError, nil, nil)
	})
	t.Run("eof", func(t *testing.T) {
		f(t, "", io.EOF, &struct{}{}, &struct{}{})
	})
	t.Run("null", func(t *testing.T) {
		f(t, "null", nil, &struct{}{}, &struct{}{})
	})
	t.Run("invalid first byte", func(t *testing.T) {
		f(t, "+", UnexpectedByteError{}, &struct{}{}, &struct{}{})
	})
	debug.FreeOSMemory()
}

func TestValDecoder_Native_Struct_Mapping(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		checkDecodeWithStandard(t, DefaultDecoder, data, ex, p1, p2)
	}
	t.Run("unexported field", func(t *testing.T) {
		f(t, ` { "a" : "abc" } `, nil, &struct {
			a string
		}{}, &struct {
			a string
		}{})
	})
	t.Run("unexported field 2", func(t *testing.T) {
		f(t, ` { "a" : "abc" } `, nil, &struct {
			a string
			B int
		}{}, &struct {
			a string
			B int
		}{})
	})
	t.Run("tag ignored 1", func(t *testing.T) {
		f(t, ` { "A" : "abc" } `, nil, &struct {
			A string `json:"-"`
		}{A: "test"}, &struct {
			A string `json:"-"`
		}{A: "test"})
	})
	t.Run("tag", func(t *testing.T) {
		f(t, ` { "b" : "abc" } `, nil, &struct {
			A string `json:"B"`
		}{A: "test"}, &struct {
			A string `json:"B"`
		}{A: "test"})
	})
	t.Run("case insensitive", func(t *testing.T) {
		f(t, ` { "a" : "abc" } `, nil, &struct {
			A string
		}{A: "test"}, &struct {
			A string
		}{A: "test"})
	})
	t.Run("case insensitive 2", func(t *testing.T) {
		f(t, ` { "A" : "abc" } `, nil, &struct {
			A string `json:"a"`
		}{A: "test"}, &struct {
			A string `json:"a"`
		}{A: "test"})
	})
	debug.FreeOSMemory()
}

func TestValDecoder_Native_Struct(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		checkDecodeWithStandard(t, DefaultDecoder, data, ex, p1, p2)
	}
	t.Run("nil receiver", func(t *testing.T) {
		f(t, "null", NilPointerReceiverError, nil, nil)
	})
	t.Run("eof", func(t *testing.T) {
		f(t, "", io.EOF, &struct {
			A string
		}{A: "test"}, &struct {
			A string
		}{A: "test"})
	})
	t.Run("invalid first byte", func(t *testing.T) {
		f(t, "+", UnexpectedByteError{}, &struct {
			A string
		}{A: "test"}, &struct {
			A string
		}{A: "test"})
	})
	t.Run("invalid null", func(t *testing.T) {
		f(t, "nul", io.EOF, &struct {
			A string
		}{A: "test"}, &struct {
			A string
		}{A: "test"})
	})
	t.Run("null", func(t *testing.T) {
		f(t, "null", nil, &struct {
			A string
		}{A: "test"}, &struct {
			A string
		}{A: "test"})
	})
	t.Run("eof after bracket", func(t *testing.T) {
		f(t, "{", io.EOF, &struct {
			A string
		}{A: "test"}, &struct {
			A string
		}{A: "test"})
	})
	t.Run("empty object", func(t *testing.T) {
		f(t, " { } ", nil, &struct {
			A string
		}{A: "test"}, &struct {
			A string
		}{A: "test"})
	})
	t.Run("invalid char after bracket", func(t *testing.T) {
		f(t, " { { ", UnexpectedByteError{}, &struct {
			A string
		}{A: "test"}, &struct {
			A string
		}{A: "test"})
	})
	t.Run("invalid field", func(t *testing.T) {
		f(t, ` { " `, io.EOF, &struct {
			A string
		}{A: "test"}, &struct {
			A string
		}{A: "test"})
	})
	t.Run("invalid field type", func(t *testing.T) {
		f(t, ` { "A" : 1 } `, UnexpectedByteError{}, &struct {
			A string
		}{A: "test"}, &struct {
			A string
		}{A: "test"})
	})
	t.Run("field skip error", func(t *testing.T) {
		f(t, ` { "b" : } `, UnexpectedByteError{}, &struct {
			A string
		}{A: "test"}, &struct {
			A string
		}{A: "test"})
	})
	t.Run("more field eof", func(t *testing.T) {
		f(t, ` { "b" : 1 `, io.EOF, &struct {
			A string
		}{A: "test"}, &struct {
			A string
		}{A: "test"})
	})
	t.Run("non empty", func(t *testing.T) {
		f(t, ` { "a" : "abc" } `, nil, &struct {
			A string
		}{A: "test"}, &struct {
			A string
		}{A: "test"})
	})
	t.Run("more field eof after comma", func(t *testing.T) {
		f(t, ` { "b" : 1 , `, io.EOF, &struct {
			A string
		}{A: "test"}, &struct {
			A string
		}{A: "test"})
	})
	t.Run("more field invalid byte after comma", func(t *testing.T) {
		f(t, ` { "b" : 1 , } `, UnexpectedByteError{}, &struct {
			A string
		}{A: "test"}, &struct {
			A string
		}{A: "test"})
	})
	t.Run("more field invalid comma", func(t *testing.T) {
		f(t, ` { "b" : 1 { `, UnexpectedByteError{}, &struct {
			A string
		}{A: "test"}, &struct {
			A string
		}{A: "test"})
	})
	t.Run("two fields", func(t *testing.T) {
		f(t, ` { "b" : 1 , "a" : "abc" } `, nil, &struct {
			A string
		}{A: "test"}, &struct {
			A string
		}{A: "test"})
	})
	debug.FreeOSMemory()
}

type testStruct struct {
	A *testStruct `json:"a"`
	C *int        `json:"c"`
	B *testStruct `json:"b"`
}

func TestValDecoder_Native_Struct_Nested(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		checkDecodeWithStandard(t, DefaultDecoder, data, ex, p1, p2)
	}
	t.Run("nested", func(t *testing.T) {
		f(t, `{"a":null,"c":1,"b":{}}`, nil, &testStruct{
			A: &testStruct{},
		}, &testStruct{
			A: &testStruct{},
		})
	})
	debug.FreeOSMemory()
}

func TestValDecoder_Native_Struct_CustomTag(t *testing.T) {
	decoder := NewDecoder(&DecoderOption{
		Tag: "jzon",
	})
	var p struct {
		A string `jzon:"b"`
	}
	err := decoder.Unmarshal([]byte(` { "b" : "c" }`), &p)
	require.NoError(t, err)
	require.Equal(t, "c", p.A)
}

func TestValDecoder_Native_Struct_Embedded_Unexported(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		checkDecodeWithStandard(t, DefaultDecoder, data, ex, p1, p2)
	}
	t.Run("not embedded", func(t *testing.T) {
		type inner struct{}
		type outer struct {
			inner inner
		}
		f(t, `{"inner":{}}`, nil, &outer{}, &outer{})
	})
	t.Run("non struct", func(t *testing.T) {
		type inner int
		type outer struct {
			inner
		}
		f(t, `{"inner":1}`, nil, &outer{}, &outer{})
	})
	t.Run("nil pointer receiver (duplicate field)", func(t *testing.T) {
		type inner struct {
			A int `json:"a"`
		}
		type inner2 inner
		type outer struct {
			*inner
			*inner2
		}
		f(t, `{"a":1}`, nil, &outer{}, &outer{})
	})
	t.Run("nil pointer receiver (field not matched)", func(t *testing.T) {
		type inner struct {
			B int
		}
		type outer struct {
			*inner
		}
		f(t, `{"a":1}`, nil, &outer{}, &outer{})
	})
	t.Run("nil pointer receiver (field matched)", func(t *testing.T) {
		type inner struct {
			A int
		}
		type outer struct {
			*inner
		}
		f(t, `{"a":1}`, NilEmbeddedPointerError, &outer{}, &outer{})
	})
}

func TestValDecoder_Native_Struct_DisallowUnknownFields(t *testing.T) {
	dec := NewDecoder(&DecoderOption{
		DisallowUnknownFields: true,
	})
	t.Run("zero field (eof)", func(t *testing.T) {
		var st struct{}
		err := dec.Unmarshal([]byte(""), &st)
		require.Equal(t, io.EOF, err)
	})
	t.Run("zero field (invalid)", func(t *testing.T) {
		var st struct{}
		err := dec.Unmarshal([]byte("["), &st)
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("zero field (null)", func(t *testing.T) {
		var st struct{}
		err := dec.Unmarshal([]byte("null"), &st)
		require.NoError(t, err)
	})
	t.Run("zero field (eof after bracket)", func(t *testing.T) {
		var st struct{}
		err := dec.Unmarshal([]byte("{"), &st)
		require.Equal(t, io.EOF, err)
	})
	t.Run("zero field (non empty)", func(t *testing.T) {
		var st struct{}
		err := dec.Unmarshal([]byte(`{"a":1}`), &st)
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("zero field (empty)", func(t *testing.T) {
		var st struct{}
		err := dec.Unmarshal([]byte(` { } `), &st)
		require.NoError(t, err)
	})
	t.Run("one field (empty)", func(t *testing.T) {
		var st struct {
			A int
		}
		err := dec.Unmarshal([]byte(` { } `), &st)
		require.NoError(t, err)
	})
	t.Run("one field (empty)", func(t *testing.T) {
		var st struct {
			A int
		}
		err := dec.Unmarshal([]byte(` { "b" : 1 } `), &st)
		require.IsType(t, UnknownFieldError(""), err)
	})
}

func TestValDecoder_Native_Struct_Quoted(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		checkDecodeWithStandard(t, DefaultDecoder, data, ex, p1, p2)
	}
	t.Run("int8", func(t *testing.T) {
		t.Run("leading space", func(t *testing.T) {
			type st struct {
				A int8 `json:",string"`
			}
			f(t, `{"a":" 1"}`, InvalidDigitError{}, &st{}, &st{})
		})
		t.Run("trailing space", func(t *testing.T) {
			type st struct {
				A int8 `json:",string"`
			}
			f(t, `{"a":"1 "}`, UnexpectedByteError{}, &st{}, &st{})
		})
		t.Run("ok", func(t *testing.T) {
			type st struct {
				A int8 `json:",string"`
			}
			f(t, `{"a":"1"}`, nil, &st{}, &st{})
		})
	})
}
