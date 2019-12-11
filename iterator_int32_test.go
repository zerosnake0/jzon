package jzon

import (
	"fmt"
	"io"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Int_ReadUint32(t *testing.T) {
	t.Run("leading space", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" 1"))
		u, err := it.ReadUint32()
		require.NoError(t, err)
		require.Equal(t, uint32(1), u)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		it := NewIterator()
		_, err := it.ReadUint32()
		require.Equal(t, io.EOF, err)
	})
	t.Run("zero", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("0"))
		i, err := it.ReadUint32()
		require.NoError(t, err)
		require.Equal(t, uint32(0), i)
	})
	t.Run("invalid first digit", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("-"))
		_, err := it.ReadUint32()
		require.IsType(t, InvalidDigitError{}, err)
	})
	t.Run("early return", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("1a"))
		u, err := it.ReadUint32()
		require.NoError(t, err)
		require.Equal(t, uint32(1), u)
	})
	t.Run("overflow 1", func(t *testing.T) {
		it := NewIterator()
		s := fmt.Sprintf("%d0", math.MaxUint32/10+1)
		it.ResetBytes([]byte(s))
		_, err := it.ReadUint32()
		require.IsType(t, IntOverflowError{}, err)
	})
	t.Run("overflow 2", func(t *testing.T) {
		it := NewIterator()
		d := math.MaxUint32 / 10
		m := math.MaxUint32 - d*10
		s := fmt.Sprintf("%d%d", d, m+1)
		it.ResetBytes([]byte(s))
		_, err := it.ReadUint32()
		require.IsType(t, IntOverflowError{}, err)
	})
	t.Run("max uint32", func(t *testing.T) {
		it := NewIterator()
		var m uint32 = math.MaxUint32
		it.ResetBytes([]byte(fmt.Sprint(m)))
		u, err := it.ReadUint32()
		require.NoError(t, err)
		require.Equal(t, m, u)
	})
	t.Run("reader", func(t *testing.T) {
		it := NewIterator()
		var m uint32 = math.MaxUint32
		it.Reset(&oneByteReader{
			b: []byte(fmt.Sprint(m)),
		})
		u, err := it.ReadUint32()
		require.NoError(t, err)
		require.Equal(t, m, u)
	})
}

func TestIterator_Int_ReadInt32(t *testing.T) {
	t.Run("leading space", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" -1"))
		i, err := it.ReadInt32()
		require.NoError(t, err)
		require.Equal(t, int32(-1), i)
	})
	t.Run("first byte error", func(t *testing.T) {
		it := NewIterator()
		_, err := it.ReadInt32()
		require.Equal(t, io.EOF, err)
	})
	t.Run("negative error", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("-"))
		_, err := it.ReadInt32()
		require.Equal(t, io.EOF, err)
	})
	t.Run("negative readUint32 error", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("-a"))
		_, err := it.ReadInt32()
		require.IsType(t, InvalidDigitError{}, err)
	})
	t.Run("negative overflow", func(t *testing.T) {
		it := NewIterator()
		s := fmt.Sprint(-math.MaxInt32 - 2)
		it.ResetBytes([]byte(s))
		_, err := it.ReadInt32()
		e, ok := err.(IntOverflowError)
		require.True(t, ok)
		require.Equal(t, "int32", e.typ)
		require.Equal(t, s, e.value)
	})
	t.Run("negative max", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(fmt.Sprint(-math.MaxInt32 - 1)))
		i, err := it.ReadInt32()
		require.NoError(t, err)
		require.Equal(t, int32(-math.MaxInt32-1), i)
	})
	t.Run("positive readUint32 error", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("a"))
		_, err := it.ReadInt32()
		require.IsType(t, InvalidDigitError{}, err)
	})
	t.Run("positive overflow", func(t *testing.T) {
		it := NewIterator()
		s := fmt.Sprint(math.MaxInt32 + 1)
		it.ResetBytes([]byte(s))
		_, err := it.ReadInt32()
		e, ok := err.(IntOverflowError)
		require.True(t, ok)
		require.Equal(t, "int32", e.typ)
		require.Equal(t, s, e.value)
	})
	t.Run("positive max", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(fmt.Sprint(math.MaxInt32)))
		i, err := it.ReadInt32()
		require.NoError(t, err)
		require.Equal(t, int32(math.MaxInt32), i)
	})
}
