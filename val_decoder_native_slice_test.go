package jzon

import (
	"io"
	"runtime"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValDecoder_Native_Slice(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
	}
	t.Run("nil pointer", func(t *testing.T) {
		f(t, "null", ErrNilPointerReceiver, nil, nil)
	})
	t.Run("eof", func(t *testing.T) {
		arr1 := []int{1, 2}
		arr2 := []int{1, 2}
		f(t, "", io.EOF, &arr1, &arr2)
	})
	t.Run("invalid null", func(t *testing.T) {
		arr1 := []int{1, 2}
		arr2 := []int{1, 2}
		f(t, "nul", io.EOF, &arr1, &arr2)
	})
	t.Run("null", func(t *testing.T) {
		arr1 := []int{1, 2}
		arr2 := []int{1, 2}
		f(t, "null", nil, &arr1, &arr2)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		arr1 := []int{1, 2}
		arr2 := []int{1, 2}
		f(t, `+`, UnexpectedByteError{}, &arr1, &arr2)
	})
	t.Run("eof after bracket", func(t *testing.T) {
		arr1 := []int{1, 2}
		arr2 := []int{1, 2}
		f(t, ` [ `, io.EOF, &arr1, &arr2)
	})
	t.Run("empty", func(t *testing.T) {
		arr1 := []int{1, 2}
		arr2 := []int{1, 2}
		f(t, ` [ ] `, nil, &arr1, &arr2)
	})
	t.Run("bad item", func(t *testing.T) {
		arr1 := []int{1, 2}
		arr2 := []int{1, 2}
		f(t, ` [ + `, InvalidDigitError{}, &arr1, &arr2)
	})
	t.Run("eof after item", func(t *testing.T) {
		arr1 := []int{1, 2}
		arr2 := []int{1, 2}
		f(t, ` [ 3 `, io.EOF, &arr1, &arr2)
	})
	t.Run("lesser item", func(t *testing.T) {
		arr1 := []int{1, 2}
		arr2 := []int{1, 2}
		f(t, ` [ 3 ] `, nil, &arr1, &arr2)
	})
	t.Run("bad item 2", func(t *testing.T) {
		arr1 := []int{1, 2}
		arr2 := []int{1, 2}
		f(t, ` [ "test" ] `, InvalidDigitError{}, &arr1, &arr2)
	})
	t.Run("invalid comma", func(t *testing.T) {
		arr1 := []int{1, 2}
		arr2 := []int{1, 2}
		f(t, ` [ 3 [ `, UnexpectedByteError{}, &arr1, &arr2)
	})
	t.Run("more item", func(t *testing.T) {
		arr1 := []int{1, 2}
		arr2 := []int{1, 2}
		f(t, ` [ 3 , 4 , 5 ] `, nil, &arr1, &arr2)
	})
	t.Run("more item error", func(t *testing.T) {
		arr1 := []int{1, 2}
		arr2 := []int{1, 2}
		f(t, ` [ 3 , 4 , "test" ] `, InvalidDigitError{}, &arr1, &arr2)
	})
	debug.FreeOSMemory()
}

func TestValDecoder_Native_Slice_Memory(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		f := func(i int) []*int {
			pi := &i
			runtime.SetFinalizer(pi, func(_ *int) {
				t.Logf("finalizing")
			})
			return []*int{pi}
		}
		arr := f(1)
		err := Unmarshal([]byte(`[]`), &arr)
		require.NoError(t, err)
		debug.FreeOSMemory()
		t.Logf("please check if the memory has been freed")
	})
	t.Run("test2", func(t *testing.T) {
		type st struct {
			p *int
		}
		f := func(i int) []*st {
			pi := &i
			runtime.SetFinalizer(pi, func(_ *int) {
				t.Logf("finalizing")
			})
			return []*st{{p: pi}}
		}
		arr := f(1)
		err := Unmarshal([]byte(`[]`), &arr)
		require.NoError(t, err)
		debug.FreeOSMemory()
		t.Logf("please check if the memory has been freed")
	})
}

func TestValDecoder_Native_Slice_AllocateError(t *testing.T) {
	type bigSt struct {
		A string
		B string
		C string
		D string
		E string
	}
	var arr []bigSt
	err := Unmarshal([]byte(`[{
		"E": "test0"
	}, {
		"E": "test1"
	}]`), &arr)
	require.NoError(t, err)
	require.Len(t, arr, 2)
	require.Equal(t, "test0", arr[0].E)
	require.Equal(t, "test1", arr[1].E)
}
