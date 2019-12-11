package jzon

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Null_ReadNull(t *testing.T) {
	t.Run("eof", func(t *testing.T) {
		it := NewIterator()
		err := it.ReadNull()
		require.Equal(t, io.EOF, err)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" a"))
		err := it.ReadNull()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("error 1", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" n"))
		err := it.ReadNull()
		require.Equal(t, io.EOF, err)
	})
	t.Run("error 2", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" na"))
		err := it.ReadNull()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("valid", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" null "))
		err := it.ReadNull()
		require.NoError(t, err)
		_, err = it.NextValueType()
		require.Equal(t, io.EOF, err)
	})
}
