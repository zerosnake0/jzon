package jzon

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Array_ReadArrayBegin(t *testing.T) {
	t.Run("leading space", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" [ ]"))
		more, err := it.ReadArrayBegin()
		require.NoError(t, err)
		require.False(t, more)
	})
	t.Run("eof", func(t *testing.T) {
		it := NewIterator()
		_, err := it.ReadArrayBegin()
		require.Equal(t, io.EOF, err)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("1"))
		_, err := it.ReadArrayBegin()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("eof after bracket", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("["))
		_, err := it.ReadArrayBegin()
		require.Equal(t, io.EOF, err)
	})
	t.Run("more", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(`["`))
		more, err := it.ReadArrayBegin()
		require.NoError(t, err)
		require.True(t, more)
	})
}

func TestIterator_Array_ReadArrayMore(t *testing.T) {
	t.Run("leading space", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" ,"))
		more, err := it.ReadArrayMore()
		require.NoError(t, err)
		require.True(t, more)
	})
	t.Run("eof", func(t *testing.T) {
		it := NewIterator()
		_, err := it.ReadArrayMore()
		require.Equal(t, io.EOF, err)
	})
	t.Run("more", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(","))
		more, err := it.ReadArrayMore()
		require.NoError(t, err)
		require.True(t, more)
	})
	t.Run("no more", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("]"))
		more, err := it.ReadArrayMore()
		require.NoError(t, err)
		require.False(t, more)
	})
	t.Run("error", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("a"))
		_, err := it.ReadArrayMore()
		require.IsType(t, UnexpectedByteError{}, err)
	})
}

func TestIterator_Array_ReadArray_Example(t *testing.T) {
	must := require.New(t)

	it := NewIterator()
	it.ResetBytes([]byte(" [ 0 , 1 , 2 ] "))
	more, err := it.ReadArrayBegin()
	must.NoError(err)
	i := 0
	for ; more; more, err = it.ReadArrayMore() {
		ri, err := it.ReadInt()
		must.NoError(err)
		must.Equal(i, ri)
		i += 1
	}
	must.NoError(err)
	_, err = it.NextValueType()
	must.Equal(io.EOF, err)
}

func TestIterator_Array_ReadArrayCB(t *testing.T) {
	t.Run("leading space", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" []"))
		err := it.ReadArrayCB(nil)
		require.NoError(t, err)
	})
	t.Run("eof", func(t *testing.T) {
		it := NewIterator()
		err := it.ReadArrayCB(nil)
		require.Equal(t, io.EOF, err)
	})
	t.Run("invalid bracket", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("a"))
		err := it.ReadArrayCB(nil)
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("eof after bracket", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("["))
		err := it.ReadArrayCB(nil)
		require.Equal(t, io.EOF, err)
	})
	t.Run("no element", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("[]"))
		err := it.ReadArrayCB(nil)
		require.NoError(t, err)
	})
	t.Run("no element", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("[]"))
		err := it.ReadArrayCB(nil)
		require.NoError(t, err)
	})
	t.Run("callback error", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("[1]"))
		e := errors.New("test")
		err := it.ReadArrayCB(func(*Iterator) error {
			return e
		})
		require.Equal(t, e, err)
	})
	t.Run("error on more", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("[1"))
		err := it.ReadArrayCB(func(it *Iterator) (err error) {
			_, err = it.ReadInt()
			return
		})
		require.Equal(t, io.EOF, err)
	})
	t.Run("error on more 2", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("[1a"))
		err := it.ReadArrayCB(func(it *Iterator) (err error) {
			_, err = it.ReadInt()
			return
		})
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("2 items", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("[1,2]"))
		err := it.ReadArrayCB(func(it *Iterator) (err error) {
			_, err = it.ReadInt()
			return
		})
		require.NoError(t, err)
	})
}

func TestIterator_Array_ReadArrayCB_Example(t *testing.T) {
	must := require.New(t)

	it := NewIterator()
	it.ResetBytes([]byte(" [ 0 , 1 , 2 ] "))
	i := 0
	err := it.ReadArrayCB(func(it *Iterator) (err error) {
		j, err := it.ReadInt()
		must.NoError(err)
		must.Equal(i, j)
		i += 1
		return nil
	})
	must.NoError(err)
	_, err = it.NextValueType()
	must.Equal(io.EOF, err)
}
