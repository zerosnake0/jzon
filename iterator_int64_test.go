package jzon

import (
	"fmt"
	"io"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Int_ReadUint64(t *testing.T) {
	t.Run("leading space", func(t *testing.T) {
		withIterator(" 1", func(it *Iterator) {
			u, err := it.ReadUint64()
			require.NoError(t, err)
			require.Equal(t, uint64(1), u)
		})
	})
	t.Run("invalid first byte", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			_, err := it.ReadUint64()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("zero", func(t *testing.T) {
		withIterator("0", func(it *Iterator) {
			i, err := it.ReadUint64()
			require.NoError(t, err)
			require.Equal(t, uint64(0), i)
		})
	})
	t.Run("invalid first digit", func(t *testing.T) {
		withIterator("-", func(it *Iterator) {
			_, err := it.ReadUint64()
			require.IsType(t, InvalidDigitError{}, err)
		})
	})
	t.Run("early return", func(t *testing.T) {
		withIterator("1a", func(it *Iterator) {
			u, err := it.ReadUint64()
			require.NoError(t, err)
			require.Equal(t, uint64(1), u)
		})
	})
	t.Run("overflow 1", func(t *testing.T) {
		s := fmt.Sprintf("%d0", math.MaxUint64/10+1)
		withIterator(s, func(it *Iterator) {
			_, err := it.ReadUint64()
			require.IsType(t, IntOverflowError{}, err)
		})
	})
	t.Run("overflow 2", func(t *testing.T) {
		d := uint64(math.MaxUint64) / 10
		m := uint64(math.MaxUint64) - d*10 + 1
		require.Less(t, m, uint64(10))
		s := fmt.Sprintf("%d%d", d, m)
		withIterator(s, func(it *Iterator) {
			_, err := it.ReadUint64()
			require.IsType(t, IntOverflowError{}, err)
		})
	})
	t.Run("max uint64", func(t *testing.T) {
		var m uint64 = math.MaxUint64
		s := fmt.Sprint(m)
		withIterator(s, func(it *Iterator) {
			u, err := it.ReadUint64()
			require.NoError(t, err)
			require.Equal(t, m, u)
		})
	})
	t.Run("reader", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			var m uint64 = math.MaxUint64
			it.Reset(&oneByteReader{
				b: fmt.Sprint(m),
			})
			u, err := it.ReadUint64()
			require.NoError(t, err)
			require.Equal(t, m, u)
		})
	})
}

func TestIterator_Int_ReadInt64(t *testing.T) {
	t.Run("leading space", func(t *testing.T) {
		withIterator(" -1", func(it *Iterator) {
			i, err := it.ReadInt64()
			require.NoError(t, err)
			require.Equal(t, int64(-1), i)
		})
	})
	t.Run("first byte error", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			_, err := it.ReadInt64()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("negative error", func(t *testing.T) {
		withIterator("-", func(it *Iterator) {
			_, err := it.ReadInt64()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("negative readUint64 error", func(t *testing.T) {
		withIterator("-a", func(it *Iterator) {
			_, err := it.ReadInt64()
			require.IsType(t, InvalidDigitError{}, err)
		})
	})
	t.Run("negative overflow", func(t *testing.T) {
		d := int64(math.MaxInt64) / 10
		m := int64(math.MaxInt64) - d*10 + 2
		require.Less(t, m, int64(10))
		s := fmt.Sprintf("-%d%d", d, m)
		withIterator(s, func(it *Iterator) {
			_, err := it.ReadInt64()
			e, ok := err.(IntOverflowError)
			require.True(t, ok)
			require.Equal(t, "int64", e.typ)
			require.Equal(t, s, e.value)
		})
	})
	t.Run("negative max", func(t *testing.T) {
		v := int64(-math.MaxInt64 - 1)
		s := fmt.Sprint(v)
		withIterator(s, func(it *Iterator) {
			i, err := it.ReadInt64()
			require.NoError(t, err)
			require.Equal(t, v, i)
		})
	})
	t.Run("positive readUint64 error", func(t *testing.T) {
		withIterator("a", func(it *Iterator) {
			_, err := it.ReadInt64()
			require.IsType(t, InvalidDigitError{}, err)
		})
	})
	t.Run("positive overflow", func(t *testing.T) {
		d := int64(math.MaxInt64 / 10)
		m := int64(math.MaxInt64-d*10) + 1
		require.Less(t, m, int64(10))
		s := fmt.Sprintf("%d%d", d, m)
		withIterator(s, func(it *Iterator) {
			_, err := it.ReadInt64()
			e, ok := err.(IntOverflowError)
			require.True(t, ok)
			require.Equal(t, "int64", e.typ)
			require.Equal(t, s, e.value)
		})
	})
	t.Run("positive max", func(t *testing.T) {
		s := fmt.Sprint(math.MaxInt64)
		withIterator(s, func(it *Iterator) {
			i, err := it.ReadInt64()
			require.NoError(t, err)
			require.Equal(t, int64(math.MaxInt64), i)
		})
	})
}
