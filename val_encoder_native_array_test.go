package jzon

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestValEncoder_Array_Error(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		arr := [...]json.Marshaler{testJsonMarshaler{
			data: []byte(`"test"`),
			err:  errors.New("test"),
		}}
		checkEncodeValueWithStandard(t, DefaultEncoder, arr)
	})
}

func TestValEncoder_Array_Empty(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		arr := [...]int{}
		checkEncodeValueWithStandard(t, DefaultEncoder, arr)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, DefaultEncoder, (*[0]int)(nil))
	})
	t.Run("empty pointer", func(t *testing.T) {
		arr := [...]int{}
		checkEncodeValueWithStandard(t, DefaultEncoder, &arr)
	})
}

func TestValEncoder_Array_Indirect(t *testing.T) {
	// len != 1
	t.Run("pointer", func(t *testing.T) {
		arr := [...]int{1, 2, 3}
		checkEncodeValueWithStandard(t, DefaultEncoder, &arr)
	})
	t.Run("non pointer", func(t *testing.T) {
		arr := [...]int{1, 2, 3}
		checkEncodeValueWithStandard(t, DefaultEncoder, arr)
	})
	t.Run("array of pointer", func(t *testing.T) {
		i := 1
		arr := [...]*int{(*int)(nil), &i}
		checkEncodeValueWithStandard(t, DefaultEncoder, arr)
	})
	// element is indirect
	t.Run("one element array", func(t *testing.T) {
		arr := [...]int{1}
		checkEncodeValueWithStandard(t, DefaultEncoder, arr)
	})
}

func TestValEncoder_Array_Direct(t *testing.T) {
	t.Run("nil element", func(t *testing.T) {
		arr := [...]*int{(*int)(nil)}
		checkEncodeValueWithStandard(t, DefaultEncoder, arr)
	})
	t.Run("non nil element", func(t *testing.T) {
		i := 1
		arr := [...]*int{&i}
		checkEncodeValueWithStandard(t, DefaultEncoder, arr)
	})
	t.Run("pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, DefaultEncoder, (*[1]*int)(nil))
	})
}
