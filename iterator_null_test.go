package jzon

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Null_ReadNull(t *testing.T) {
	t.Run("eof", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			err := it.ReadNull()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("invalid first byte", func(t *testing.T) {
		withIterator(" a", func(it *Iterator) {
			err := it.ReadNull()
			require.IsType(t, UnexpectedByteError{}, err)
		})
	})
	t.Run("error 1", func(t *testing.T) {
		withIterator(" n", func(it *Iterator) {
			err := it.ReadNull()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("error 2", func(t *testing.T) {
		withIterator(" na", func(it *Iterator) {
			err := it.ReadNull()
			require.IsType(t, UnexpectedByteError{}, err)
		})
	})
	t.Run("valid", func(t *testing.T) {
		withIterator(" null ", func(it *Iterator) {
			err := it.ReadNull()
			require.NoError(t, err)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
}
