package jzon

import (
	"testing"
)

type testIfaceImpl2 struct {
	Field string
}

func (testIfaceImpl2) Foo() {}

type testIfaceImpl3 struct {
	Field string
}

func (*testIfaceImpl3) Foo() {}

func TestValEncoder_Interface(t *testing.T) {
	t.Run("builtin", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			var i interface{}
			checkEncodeValueWithStandard(t, &i, nil)
		})
		t.Run("non nil", func(t *testing.T) {
			var i interface{} = 1
			checkEncodeValueWithStandard(t, &i, nil)
		})
	})
	t.Run("eface", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			type I interface{}
			var i I
			checkEncodeValueWithStandard(t, &i, nil)
		})
		t.Run("non nil", func(t *testing.T) {
			type I interface{}
			var i I = 1
			checkEncodeValueWithStandard(t, &i, nil)
		})
	})
	t.Run("iface", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			var i testIface
			checkEncodeValueWithStandard(t, &i, nil)
		})
		t.Run("non nil", func(t *testing.T) {
			var i testIface = testIfaceImpl2{
				Field: "test",
			}
			checkEncodeValueWithStandard(t, &i, nil)
		})
		t.Run("non nil 2", func(t *testing.T) {
			var i testIface = &testIfaceImpl2{
				Field: "test2",
			}
			checkEncodeValueWithStandard(t, &i, nil)
		})
		t.Run("pointer receiver", func(t *testing.T) {
			var i testIface = &testIfaceImpl3{
				Field: "test3",
			}
			checkEncodeValueWithStandard(t, &i, nil)
		})
	})
}

func TestValEncoder_Interface_OmitEmpty(t *testing.T) {
	t.Run("eface", func(t *testing.T) {
		type st struct {
			I interface{} `json:",omitempty"`
		}
		t.Run("nil", func(t *testing.T) {
			checkEncodeValueWithStandard(t, st{}, nil)
		})
		t.Run("non nil", func(t *testing.T) {
			checkEncodeValueWithStandard(t, st{I: true}, nil)
		})
		t.Run("nil pointer", func(t *testing.T) {
			checkEncodeValueWithStandard(t, st{
				I: (*testJSONMarshaler2)(nil),
			}, nil)
		})
	})
	t.Run("iface", func(t *testing.T) {
		type st struct {
			I testIface `json:",omitempty"`
		}
		t.Run("nil", func(t *testing.T) {
			checkEncodeValueWithStandard(t, st{}, nil)
		})
		t.Run("value receiver", func(t *testing.T) {
			t.Run("non nil", func(t *testing.T) {
				checkEncodeValueWithStandard(t, st{
					I: testIfaceImpl{},
				}, nil)
			})
			t.Run("empty struct", func(t *testing.T) {
				checkEncodeValueWithStandard(t, st{
					I: testIfaceImpl2{},
				}, nil)
			})
			t.Run("non empty struct", func(t *testing.T) {
				checkEncodeValueWithStandard(t, st{
					I: testIfaceImpl2{
						Field: "test2",
					},
				}, nil)
			})
		})
		t.Run("pointer receiver", func(t *testing.T) {
			t.Run("nil pointer", func(t *testing.T) {
				checkEncodeValueWithStandard(t, st{
					I: (*testIfaceImpl3)(nil),
				}, nil)
			})
			t.Run("empty struct", func(t *testing.T) {
				checkEncodeValueWithStandard(t, st{
					I: &testIfaceImpl3{},
				}, nil)
			})
			t.Run("non empty struct", func(t *testing.T) {
				checkEncodeValueWithStandard(t, st{
					I: &testIfaceImpl3{
						Field: "test3",
					},
				}, nil)
			})
		})
	})
}
