package jzon

import (
	"encoding/json"
	"errors"
	"runtime/debug"
	"testing"
)

func TestValEncoder_Slice_Error(t *testing.T) {
	t.Run("chain error", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			(*sliceEncoder)(nil).Encode(nil, s)
		})
	})
	t.Run("error", func(t *testing.T) {
		e := errors.New("test")
		arr := []json.Marshaler{testJsonMarshaler{
			data: []byte(`"test"`),
			err:  e,
		}}
		checkEncodeValueWithStandard(t, arr, e)
	})
	debug.FreeOSMemory()
}

func TestValEncoder_Slice_Empty(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var arr []int
		checkEncodeValueWithStandard(t, arr, nil)
	})
	t.Run("empty", func(t *testing.T) {
		arr := make([]int, 0)
		checkEncodeValueWithStandard(t, arr, nil)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, (*[]int)(nil), nil)
	})
	t.Run("empty pointer", func(t *testing.T) {
		arr := make([]int, 0)
		checkEncodeValueWithStandard(t, &arr, nil)
	})
	debug.FreeOSMemory()
}

func TestValEncoder_Slice(t *testing.T) {
	t.Run("pointer", func(t *testing.T) {
		arr := []int{1, 2, 3}
		checkEncodeValueWithStandard(t, &arr, nil)
	})
	t.Run("non pointer", func(t *testing.T) {
		arr := []int{1, 2, 3}
		checkEncodeValueWithStandard(t, arr, nil)
	})
	t.Run("slice of pointer", func(t *testing.T) {
		i := 1
		arr := []*int{(*int)(nil), &i}
		checkEncodeValueWithStandard(t, arr, nil)
	})
	debug.FreeOSMemory()
}
