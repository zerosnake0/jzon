package jzon

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestValEncoder_Array_Error(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		e := errors.New("test")
		arr := [...]json.Marshaler{testJsonMarshaler{
			data: `"test"`,
			err:  e,
		}}
		checkEncodeValueWithStandard(t, arr, e)
	})
}

func TestValEncoder_Array_Empty(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		arr := [...]int{}
		checkEncodeValueWithStandard(t, arr, nil)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, (*[0]int)(nil), nil)
	})
	t.Run("empty pointer", func(t *testing.T) {
		arr := [...]int{}
		checkEncodeValueWithStandard(t, &arr, nil)
	})
}

func TestValEncoder_Array_Indirect(t *testing.T) {
	// len != 1
	t.Run("pointer", func(t *testing.T) {
		arr := [...]int{1, 2, 3}
		checkEncodeValueWithStandard(t, &arr, nil)
	})
	t.Run("non pointer", func(t *testing.T) {
		arr := [...]int{1, 2, 3}
		checkEncodeValueWithStandard(t, arr, nil)
	})
	t.Run("array of pointer", func(t *testing.T) {
		i := 1
		arr := [...]*int{(*int)(nil), &i}
		checkEncodeValueWithStandard(t, arr, nil)
	})
	// element is indirect
	t.Run("one element array", func(t *testing.T) {
		arr := [...]int{1}
		checkEncodeValueWithStandard(t, arr, nil)
	})
}

func TestValEncoder_Array_Direct(t *testing.T) {
	t.Run("nil element", func(t *testing.T) {
		arr := [...]*int{(*int)(nil)}
		checkEncodeValueWithStandard(t, arr, nil)
	})
	t.Run("non nil element", func(t *testing.T) {
		i := 1
		arr := [...]*int{&i}
		checkEncodeValueWithStandard(t, arr, nil)
	})
	t.Run("pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, (*[1]*int)(nil), nil)
	})
}
