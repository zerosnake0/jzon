package jzon

import (
	"io"
	"runtime"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValDecoder_Native_Array(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		checkStandard(t, DefaultDecoder, data, ex, p1, p2)
	}
	t.Run("nil pointer", func(t *testing.T) {
		f(t, "null", NilPointerReceiverError, nil, nil)
	})
	t.Run("eof", func(t *testing.T) {
		arr1 := [...]int{1, 2}
		arr2 := [...]int{1, 2}
		f(t, "", io.EOF, &arr1, &arr2)
	})
	t.Run("null", func(t *testing.T) {
		arr1 := [...]int{1, 2}
		arr2 := [...]int{1, 2}
		f(t, "null", nil, &arr1, &arr2)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		arr1 := [...]int{1, 2}
		arr2 := [...]int{1, 2}
		f(t, "+", UnexpectedByteError{}, &arr1, &arr2)
	})
	t.Run("eof after bracket", func(t *testing.T) {
		arr1 := [...]int{1, 2}
		arr2 := [...]int{1, 2}
		f(t, "[", io.EOF, &arr1, &arr2)
	})
	t.Run("empty", func(t *testing.T) {
		arr1 := [...]int{1, 2}
		arr2 := [...]int{1, 2}
		f(t, " [ ] ", nil, &arr1, &arr2)
	})
	t.Run("element error", func(t *testing.T) {
		arr1 := [...]int{1, 2}
		arr2 := [...]int{1, 2}
		f(t, ` [ " `, InvalidDigitError{}, &arr1, &arr2)
	})
	t.Run("null element", func(t *testing.T) {
		arr1 := [...]int{1, 2}
		arr2 := [...]int{1, 2}
		f(t, ` [ null ] `, nil, &arr1, &arr2)
	})
	t.Run("eof after element", func(t *testing.T) {
		arr1 := [...]int{1, 2}
		arr2 := [...]int{1, 2}
		f(t, ` [ 2 `, io.EOF, &arr1, &arr2)
	})
	t.Run("lesser element", func(t *testing.T) {
		arr1 := [...]int{1, 2}
		arr2 := [...]int{1, 2}
		f(t, ` [ 2 ] `, nil, &arr1, &arr2)
	})
	t.Run("invalid comma", func(t *testing.T) {
		arr1 := [...]int{1, 2}
		arr2 := [...]int{1, 2}
		f(t, ` [ 2 [ `, UnexpectedByteError{}, &arr1, &arr2)
	})
	t.Run("more element error", func(t *testing.T) {
		arr1 := [...]int{1, 2}
		arr2 := [...]int{1, 2}
		f(t, ` [ 2 , 3 , `, io.EOF, &arr1, &arr2)
	})
	t.Run("more element", func(t *testing.T) {
		arr1 := [...]int{1, 2}
		arr2 := [...]int{1, 2}
		f(t, ` [ 2 , 3 , "test"]`, nil, &arr1, &arr2)
	})
	debug.FreeOSMemory()
}

func TestValDecoder_Native_Array_Memory(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		f := func(i int) [1]*int {
			pi := &i
			runtime.SetFinalizer(pi, func(_ *int) {
				t.Logf("finalizing")
			})
			return [1]*int{pi}
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
		f := func(i int) [1]*st {
			pi := &i
			runtime.SetFinalizer(pi, func(_ *int) {
				t.Logf("finalizing")
			})
			return [1]*st{{p: pi}}
		}
		arr := f(1)
		err := Unmarshal([]byte(`[]`), &arr)
		require.NoError(t, err)
		debug.FreeOSMemory()
		t.Logf("please check if the memory has been freed")
	})
}
