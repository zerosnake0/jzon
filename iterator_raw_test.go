package jzon

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Raw_ReadRaw(t *testing.T) {
	t.Run("eof", func(t *testing.T) {
		it := NewIterator()
		_, err := it.ReadRaw()
		require.Equal(t, io.EOF, err)
	})
	t.Run("error", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` + `))
		_, err := it.ReadRaw()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("no error", func(t *testing.T) {
		it := NewIterator()
		exp := `{ " " : null }`

		data := []byte(" " + exp + " ")
		it.ResetBytes(data)
		raw, err := it.ReadRaw()
		require.NoError(t, err)
		require.Equal(t, exp, string(raw))

		copy(data, []byte(exp+"  ")) // modify content
		require.Equal(t, exp, string(raw))

		_, err = it.ReadRaw()
		require.Equal(t, io.EOF, err)
	})
}

func TestIterator_Raw_AppendRaw(t *testing.T) {
	it := NewIterator()
	data := []byte(`{}`)
	it.ResetBytes(data)
	prefix := []byte(`test`)
	b, err := it.AppendRaw(prefix)
	require.NoError(t, err)
	require.Equal(t, append(prefix, data...), b)
}
