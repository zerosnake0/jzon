package jzon

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Bool_ReadBool(t *testing.T) {
	t.Run("eof", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			_, err := it.ReadBool()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("invalid fist byte", func(t *testing.T) {
		withIterator(" a", func(it *Iterator) {
			_, err := it.ReadBool()
			require.IsType(t, UnexpectedByteError{}, err)
		})
	})
	t.Run("true", func(t *testing.T) {
		withIterator(" true ", func(it *Iterator) {
			b, err := it.ReadBool()
			require.NoError(t, err)
			require.True(t, b)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("false", func(t *testing.T) {
		withIterator(" false ", func(it *Iterator) {
			b, err := it.ReadBool()
			require.NoError(t, err)
			require.False(t, b)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
}
