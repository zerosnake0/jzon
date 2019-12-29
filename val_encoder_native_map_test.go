package jzon

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"testing"
)

func TestValEncoder_Map_Error(t *testing.T) {
	t.Run("chain error", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			(*directMapEncoder)(nil).Encode(nil, s, nil)
		})
	})
	t.Run("element error", func(t *testing.T) {
		e := errors.New("test")
		checkEncodeValueWithStandard(t, map[string]json.Marshaler{
			"key": testJsonMarshaler{
				err: e,
			},
		}, e)
	})
}

func TestValEncoder_Map(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var m map[string]int
		checkEncodeValueWithStandard(t, m, nil)
	})
	t.Run("nil pointer", func(t *testing.T) {
		var m map[string]int
		checkEncodeValueWithStandard(t, &m, nil)
	})
	t.Run("simple", func(t *testing.T) {
		m := map[string]int{"1": 2}
		checkEncodeValueWithStandard(t, m, nil)
	})
	t.Run("pointer elem", func(t *testing.T) {
		i := 3
		m := map[string]*int{
			"1": (*int)(nil),
			"2": &i,
		}
		checkEncodeValueWithStandard(t, m, nil)
	})
}

type testMapIntKey2 int

func (i testMapIntKey2) MarshalText() ([]byte, error) {
	return []byte(strconv.Itoa(int(i * 2))), nil
}

type testMapStringKey2 string

func (s testMapStringKey2) MarshalText() ([]byte, error) {
	return []byte(strconv.Itoa(len(s))), nil
}

func TestValEncoder_Native_Map_KeyEncoder_TextMarshaler(t *testing.T) {
	t.Run("not supported", func(t *testing.T) {
		type key testTextMarshaler
		m := map[key]int{{
			data: "a",
		}: 1}
		checkEncodeValueWithStandard(t, m, TypeNotSupportedError(""))
	})
	t.Run("marshaler 1-non pointer", func(t *testing.T) {
		type key = testTextMarshaler
		m := map[key]int{{
			data: "a",
		}: 1}
		checkEncodeValueWithStandard(t, m, nil)
	})
	t.Run("marshaler 1-non pointer (err)", func(t *testing.T) {
		type key = testTextMarshaler
		e := errors.New("test")
		m := map[key]int{{
			data: "a",
			err:  e,
		}: 1}
		checkEncodeValueWithStandard(t, m, e)
	})
	t.Run("marshaler 1-pointer", func(t *testing.T) {
		type key = *testTextMarshaler
		m := map[key]int{{
			data: "a",
		}: 1}
		checkEncodeValueWithStandard(t, m, nil)
	})
	t.Run("marshaler 1-pointer (err)", func(t *testing.T) {
		e := errors.New("test")
		type key = *testTextMarshaler
		m := map[key]int{{
			data: "a",
			err:  e,
		}: 1}
		checkEncodeValueWithStandard(t, m, e)
	})
	t.Run("marshaler 1-pointer (nil)", func(t *testing.T) {
		type key = *testTextMarshaler
		m := map[key]int{nil: 1}
		checkEncodeValueWithStandard(t, m, runtimeErrorType)
	})
	t.Run("marshaler 2-non pointer", func(t *testing.T) {
		type key = testTextMarshaler2
		m := map[key]int{{
			data: "a",
		}: 1}
		checkEncodeValueWithStandard(t, m, TypeNotSupportedError(""))
	})
	t.Run("marshaler 2-pointer", func(t *testing.T) {
		type key = *testTextMarshaler2
		m := map[key]int{{
			data: "a",
		}: 1}
		checkEncodeValueWithStandard(t, m, nil)
	})
	t.Run("marshaler 2-pointer (err)", func(t *testing.T) {
		e := errors.New("test")
		type key = *testTextMarshaler2
		m := map[key]int{{
			data: "a",
			err:  e,
		}: 1}
		checkEncodeValueWithStandard(t, m, e)
	})
	t.Run("marshaler 2-pointer (nil)", func(t *testing.T) {
		type key = *testTextMarshaler2
		m := map[key]int{nil: 1}
		checkEncodeValueWithStandard(t, m, runtimeErrorType)
	})
	t.Run("int key", func(t *testing.T) {
		m := map[testMapIntKey2]testMapIntKey2{
			1: 2,
		}
		checkEncodeValueWithStandard(t, m, nil)
	})
	t.Run("string key", func(t *testing.T) {
		// the MarshalText of the key is ignored
		m := map[testMapStringKey2]testMapStringKey2{
			"key": "value",
		}
		checkEncodeValueWithStandard(t, m, nil)
	})
}

func TestValEncoder_Native_Map_KeyEncoder_String(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		type key string
		m := map[key]int{"a": 1}
		checkEncodeValueWithStandard(t, m, nil)
	})
}

func TestValEncoder_Native_Map_KeyEncoder_Int(t *testing.T) {
	t.Run("int8", func(t *testing.T) {
		type key int8
		m := map[key]key{
			math.MaxInt8: math.MaxInt8,
			math.MinInt8: math.MinInt8,
		}
		checkEncodeValueWithStandard(t, m, nil)
	})
	t.Run("int16", func(t *testing.T) {
		type key int16
		m := map[key]key{
			math.MaxInt16: math.MaxInt16,
			math.MinInt16: math.MinInt16,
		}
		checkEncodeValueWithStandard(t, m, nil)
	})
	t.Run("int32", func(t *testing.T) {
		type key int32
		m := map[key]key{
			math.MaxInt32: math.MaxInt32,
			math.MinInt32: math.MinInt32,
		}
		checkEncodeValueWithStandard(t, m, nil)
	})
	t.Run("int64", func(t *testing.T) {
		type key int64
		m := map[key]key{
			math.MaxInt64: math.MaxInt64,
			math.MinInt64: math.MinInt64,
		}
		checkEncodeValueWithStandard(t, m, nil)
	})
}

func TestValEncoder_Native_Map_KeyEncoder_Uint(t *testing.T) {
	t.Run("uint8", func(t *testing.T) {
		type key uint8
		m := map[key]key{
			math.MaxUint8: math.MaxUint8,
		}
		checkEncodeValueWithStandard(t, m, nil)
	})
	t.Run("uint16", func(t *testing.T) {
		type key uint16
		m := map[key]key{
			math.MaxUint16: math.MaxUint16,
		}
		checkEncodeValueWithStandard(t, m, nil)
	})
	t.Run("uint32", func(t *testing.T) {
		type key uint32
		m := map[key]key{
			math.MaxUint32: math.MaxUint32,
		}
		checkEncodeValueWithStandard(t, m, nil)
	})
	t.Run("uint64", func(t *testing.T) {
		type key uint64
		m := map[key]key{
			math.MaxUint64: math.MaxUint64,
		}
		checkEncodeValueWithStandard(t, m, nil)
	})
}
