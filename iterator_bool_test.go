package jzon

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Bool_ReadBool(t *testing.T) {
	t.Run("eof", func(t *testing.T) {
		it := NewIterator()
		_, err := it.ReadBool()
		require.Equal(t, io.EOF, err)
	})
	t.Run("invalid fist byte", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" a"))
		_, err := it.ReadBool()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("true", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" true "))
		b, err := it.ReadBool()
		require.NoError(t, err)
		require.True(t, b)
		_, err = it.NextValueType()
		require.Equal(t, io.EOF, err)
	})
	t.Run("false", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" false "))
		b, err := it.ReadBool()
		require.NoError(t, err)
		require.False(t, b)
		_, err = it.NextValueType()
		require.Equal(t, io.EOF, err)
	})
}
