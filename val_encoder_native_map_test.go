package jzon

import "testing"

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

func TestValEncoder_Native_Map_KeyEncoder_TextMarshaler(t *testing.T) {
	t.Run("", func(t *testing.T) {
		type key testTextMarshaler
		m := map[key]int{{
			data: "a",
		}: 1}
		checkEncodeValueWithStandard(t, DefaultEncoder, m,
			TypeNotSupportedError(""))
	})
}

func TestValEncoder_Native_Map_KeyEncoder_String(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		type key string
		m := map[key]int{"a": 1}
		checkEncodeValueWithStandard(t, DefaultEncoder, m, nil)
	})
}
