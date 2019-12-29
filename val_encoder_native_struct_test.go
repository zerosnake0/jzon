package jzon

import (
	"testing"

	"github.com/stretchr/testify/require"
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
	t.Run("tag ignored 1", func(t *testing.T) {
		checkEncodeValueWithStandard(t, DefaultEncoder, &struct {
			A string `json:"-"`
		}{
			A: "test",
		}, nil)
	})
	t.Run("tag", func(t *testing.T) {
		checkEncodeValueWithStandard(t, DefaultEncoder, &struct {
			A string `json:"B"`
		}{
			A: "test",
		}, nil)
	})
}

func TestValEncoder_Native_Struct_CustomTag(t *testing.T) {
	encoder := NewEncoder(&EncoderOption{
		Tag: "jzon",
	})
	b, err := encoder.Marshal(&struct {
		A string `jzon:"b"`
	}{
		A: "test",
	})
	require.NoError(t, err)
	require.Equal(t, `{"b":"test"}`, string(b))
}

func TestValEncoder_Native_Struct_Embedded_Unexported(t *testing.T) {
	t.Run("not embedded", func(t *testing.T) {
		type inner struct{}
		type outer struct {
			inner inner
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, &outer{
			inner: inner{},
		}, nil)
	})
	t.Run("non struct", func(t *testing.T) {
		type inner int
		type outer struct {
			inner
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, &outer{
			inner: 1,
		}, nil)
	})
	t.Run("duplicate field", func(t *testing.T) {
		type inner struct {
			A int `json:"a"`
		}
		type inner2 inner
		type outer struct {
			*inner
			*inner2
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, &outer{}, nil)
		checkEncodeValueWithStandard(t, DefaultEncoder, &outer{
			inner: &inner{A: 1},
		}, nil)
		checkEncodeValueWithStandard(t, DefaultEncoder, &outer{
			inner:  &inner{A: 1},
			inner2: &inner2{A: 2},
		}, nil)
	})
}
