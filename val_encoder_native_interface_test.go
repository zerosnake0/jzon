package jzon

import (
	"testing"
)

type testIfaceImpl2 struct {
	Field string
}

func (testIfaceImpl2) Foo() {
}

func TestValEncoder_Interface(t *testing.T) {
	t.Run("builtin (nil)", func(t *testing.T) {
		var i interface{}
		checkEncodeValueWithStandard(t, DefaultEncoder, &i, nil)
	})
	t.Run("builtin", func(t *testing.T) {
		var i interface{} = 1
		checkEncodeValueWithStandard(t, DefaultEncoder, &i, nil)
	})
	t.Run("eface (nil)", func(t *testing.T) {
		type I interface{}
		var i I
		checkEncodeValueWithStandard(t, DefaultEncoder, &i, nil)
	})
	t.Run("eface", func(t *testing.T) {
		type I interface{}
		var i I = 1
		checkEncodeValueWithStandard(t, DefaultEncoder, &i, nil)
	})
	t.Run("iface (nil)", func(t *testing.T) {
		var i testIface
		checkEncodeValueWithStandard(t, DefaultEncoder, &i, nil)
	})
	t.Run("iface", func(t *testing.T) {
		var i testIface = testIfaceImpl2{
			Field: "test",
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, &i, nil)
	})
	t.Run("iface 2", func(t *testing.T) {
		var i testIface = &testIfaceImpl2{
			Field: "test2",
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, &i, nil)
	})
}
