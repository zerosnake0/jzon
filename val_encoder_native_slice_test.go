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
			(*sliceEncoder)(nil).Encode(nil, s, nil)
		})
	})
	t.Run("error", func(t *testing.T) {
		e := errors.New("test")
		arr := []json.Marshaler{testJSONMarshaler{
			data: `"test"`,
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

func TestValEncoder_Slice_OmitEmpty(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		type S []int
		type st struct {
			S S `json:",omitempty"`
		}
		t.Run("nil", func(t *testing.T) {
			checkEncodeValueWithStandard(t, st{}, nil)
		})
		t.Run("empty", func(t *testing.T) {
			checkEncodeValueWithStandard(t, st{
				S: S{},
			}, nil)
		})
		t.Run("non empty", func(t *testing.T) {
			checkEncodeValueWithStandard(t, st{
				S: S{1},
			}, nil)
		})
	})
}
