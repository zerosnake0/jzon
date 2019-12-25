package jzon

import (
	"strconv"
	"testing"
)

func TestValEncoder_Map(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var m map[string]int
		checkEncodeValueWithStandard(t, DefaultEncoder, m, nil)
	})
	t.Run("nil pointer", func(t *testing.T) {
		var m map[string]int
		checkEncodeValueWithStandard(t, DefaultEncoder, &m, nil)
	})
	t.Run("simple", func(t *testing.T) {
		m := map[string]int{"1": 2}
		checkEncodeValueWithStandard(t, DefaultEncoder, m, nil)
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
		checkEncodeValueWithStandard(t, DefaultEncoder, m,
			TypeNotSupportedError(""))
	})
	t.Run("marshaler 1-non pointer", func(t *testing.T) {
		type key = testTextMarshaler
		m := map[key]int{{
			data: "a",
		}: 1}
		checkEncodeValueWithStandard(t, DefaultEncoder, m, nil)
	})
	t.Run("marshaler 1-pointer", func(t *testing.T) {
		type key = *testTextMarshaler
		m := map[key]int{{
			data: "a",
		}: 1}
		checkEncodeValueWithStandard(t, DefaultEncoder, m, nil)
	})
	t.Run("marshaler 2-non pointer", func(t *testing.T) {
		type key = testTextMarshaler2
		m := map[key]int{{
			data: "a",
		}: 1}
		checkEncodeValueWithStandard(t, DefaultEncoder, m, TypeNotSupportedError(""))
	})
	t.Run("marshaler 2-pointer", func(t *testing.T) {
		type key = *testTextMarshaler2
		m := map[key]int{{
			data: "a",
		}: 1}
		checkEncodeValueWithStandard(t, DefaultEncoder, m, nil)
	})
	t.Run("int key", func(t *testing.T) {
		m := map[testMapIntKey2]testMapIntKey2{
			1: 2,
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, m, nil)
	})
	t.Run("string key", func(t *testing.T) {
		// the MarshalText of the key is ignored
		m := map[testMapStringKey2]testMapStringKey2{
			"key": "value",
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, m, nil)
	})
}

func TestValEncoder_Native_Map_KeyEncoder_String(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		type key string
		m := map[key]int{"a": 1}
		checkEncodeValueWithStandard(t, DefaultEncoder, m, nil)
	})
}
