package jzon

import (
	"testing"
)

func TestValEncoder_Native_Struct_Zero_Field(t *testing.T) {
	t.Run("pointer", func(t *testing.T) {
		type a struct{}
		checkEncodeValueWithStandard(t, DefaultEncoder, (*a)(nil), nil)
	})
	t.Run("non pointer", func(t *testing.T) {
		type a struct{}
		checkEncodeValueWithStandard(t, DefaultEncoder, a{}, nil)
	})
}

func TestValEncoder_Native_Struct_Mapping(t *testing.T) {
	t.Run("unexported field", func(t *testing.T) {
		checkEncodeValueWithStandard(t, DefaultEncoder, &struct {
			a string
		}{
			a: "abc",
		}, nil)
	})
	t.Run("unexported field 2", func(t *testing.T) {
		checkEncodeValueWithStandard(t, DefaultEncoder, &struct {
			a string
			B int
		}{
			a: "abc",
			B: 1,
		}, nil)
	})
}
