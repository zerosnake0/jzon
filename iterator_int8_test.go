package jzon

import (
	"fmt"
	"io"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Int_ReadUint8(t *testing.T) {
	t.Run("leading space", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" 1"))
		u, err := it.ReadUint8()
		require.NoError(t, err)
		require.Equal(t, uint8(1), u)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		it := NewIterator()
		_, err := it.ReadUint8()
		require.Equal(t, io.EOF, err)
	})
	t.Run("zero", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("0"))
		i, err := it.ReadUint8()
		require.NoError(t, err)
		require.Equal(t, uint8(0), i)
	})
	t.Run("invalid first digit", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("-"))
		_, err := it.ReadUint8()
		require.IsType(t, InvalidDigitError{}, err)
	})
	t.Run("early return", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("1a"))
		u, err := it.ReadUint8()
		require.NoError(t, err)
		require.Equal(t, uint8(1), u)
	})
	t.Run("early return2", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("12a"))
		u, err := it.ReadUint8()
		require.NoError(t, err)
		require.Equal(t, uint8(12), u)
	})
	t.Run("overflow 1", func(t *testing.T) {
		it := NewIterator()
		s := fmt.Sprintf("%d0", math.MaxUint8/10+1)
		it.ResetBytes([]byte(s))
		_, err := it.ReadUint8()
		require.IsType(t, IntOverflowError{}, err)
	})
	t.Run("overflow 2", func(t *testing.T) {
		it := NewIterator()
		d := math.MaxUint8 / 10
		m := math.MaxUint8 - d*10
		s := fmt.Sprintf("%d%d", d, m+1)
		it.ResetBytes([]byte(s))
		_, err := it.ReadUint8()
		require.IsType(t, IntOverflowError{}, err)
	})
	t.Run("max uint8", func(t *testing.T) {
		it := NewIterator()
		var m uint8 = math.MaxUint8
		it.ResetBytes([]byte(fmt.Sprint(m)))
		u, err := it.ReadUint8()
		require.NoError(t, err)
		require.Equal(t, m, u)
	})
	t.Run("reader", func(t *testing.T) {
		it := NewIterator()
		var m uint8 = math.MaxUint8
		it.Reset(&oneByteReader{
			b: []byte(fmt.Sprint(m)),
		})
		u, err := it.ReadUint8()
		require.NoError(t, err)
		require.Equal(t, m, u)
	})
	t.Run("all values", func(t *testing.T) {
		it := NewIterator()
		for i := 0; i <= math.MaxUint8; i++ {
			it.ResetBytes([]byte(fmt.Sprint(i)))
			j, err := it.ReadUint8()
			require.NoError(t, err, "value %d", i)
			require.Equal(t, uint8(i), j)
		}
	})
}

func TestIterator_Int_ReadInt8(t *testing.T) {
	t.Run("leading space", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" -1"))
		i, err := it.ReadInt8()
		require.NoError(t, err)
		require.Equal(t, int8(-1), i)
	})
	t.Run("first byte error", func(t *testing.T) {
		it := NewIterator()
		_, err := it.ReadInt8()
		require.Equal(t, io.EOF, err)
	})
	t.Run("negative error", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("-"))
		_, err := it.ReadInt8()
		require.Equal(t, io.EOF, err)
	})
	t.Run("negative readUint8 error", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("-a"))
		_, err := it.ReadInt8()
		require.IsType(t, InvalidDigitError{}, err)
	})
	t.Run("negative overflow", func(t *testing.T) {
		it := NewIterator()
		s := fmt.Sprint(-math.MaxInt8 - 2)
		it.ResetBytes([]byte(s))
		_, err := it.ReadInt8()
		e, ok := err.(IntOverflowError)
		require.True(t, ok)
		require.Equal(t, "int8", e.typ)
		require.Equal(t, s, e.value)
	})
	t.Run("negative max", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(fmt.Sprint(-math.MaxInt8 - 1)))
		i, err := it.ReadInt8()
		require.NoError(t, err)
		require.Equal(t, int8(-math.MaxInt8-1), i)
	})
	t.Run("positive readUint8 error", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("a"))
		_, err := it.ReadInt8()
		require.IsType(t, InvalidDigitError{}, err)
	})
	t.Run("positive overflow", func(t *testing.T) {
		it := NewIterator()
		s := fmt.Sprint(math.MaxInt8 + 1)
		it.ResetBytes([]byte(s))
		_, err := it.ReadInt8()
		e, ok := err.(IntOverflowError)
		require.True(t, ok)
		require.Equal(t, "int8", e.typ)
		require.Equal(t, s, e.value)
	})
	t.Run("positive max", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(fmt.Sprint(math.MaxInt8)))
		i, err := it.ReadInt8()
		require.NoError(t, err)
		require.Equal(t, int8(math.MaxInt8), i)
	})
}
