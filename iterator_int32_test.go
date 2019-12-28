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
		withIterator(" 1", func(it *Iterator) {
			u, err := it.ReadUint32()
			require.NoError(t, err)
			require.Equal(t, uint32(1), u)
		})
	})
	t.Run("invalid first byte", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			_, err := it.ReadUint32()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("zero", func(t *testing.T) {
		withIterator("0", func(it *Iterator) {
			i, err := it.ReadUint32()
			require.NoError(t, err)
			require.Equal(t, uint32(0), i)
		})
	})
	t.Run("invalid first digit", func(t *testing.T) {
		withIterator("-", func(it *Iterator) {
			_, err := it.ReadUint32()
			require.IsType(t, InvalidDigitError{}, err)
		})
	})
	t.Run("early return", func(t *testing.T) {
		withIterator("1a", func(it *Iterator) {
			u, err := it.ReadUint32()
			require.NoError(t, err)
			require.Equal(t, uint32(1), u)
		})
	})
	t.Run("overflow 1", func(t *testing.T) {
		s := fmt.Sprintf("%d0", math.MaxUint32/10+1)
		withIterator(s, func(it *Iterator) {
			_, err := it.ReadUint32()
			require.IsType(t, IntOverflowError{}, err)
		})
	})
	t.Run("overflow 2", func(t *testing.T) {
		d := math.MaxUint32 / 10
		m := math.MaxUint32 - d*10
		s := fmt.Sprintf("%d%d", d, m+1)
		withIterator(s, func(it *Iterator) {
			_, err := it.ReadUint32()
			require.IsType(t, IntOverflowError{}, err)
		})
	})
	t.Run("max uint32", func(t *testing.T) {
		var m uint32 = math.MaxUint32
		s := fmt.Sprint(m)
		withIterator(s, func(it *Iterator) {
			u, err := it.ReadUint32()
			require.NoError(t, err)
			require.Equal(t, m, u)
		})
	})
	t.Run("reader", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			var m uint32 = math.MaxUint32
			it.Reset(&oneByteReader{
				b: fmt.Sprint(m),
			})
			u, err := it.ReadUint32()
			require.NoError(t, err)
			require.Equal(t, m, u)
		})
	})
}

func TestIterator_Int_ReadInt32(t *testing.T) {
	t.Run("leading space", func(t *testing.T) {
		withIterator(" -1", func(it *Iterator) {
			i, err := it.ReadInt32()
			require.NoError(t, err)
			require.Equal(t, int32(-1), i)
		})
	})
	t.Run("first byte error", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			_, err := it.ReadInt32()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("negative error", func(t *testing.T) {
		withIterator("-", func(it *Iterator) {
			_, err := it.ReadInt32()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("negative readUint32 error", func(t *testing.T) {
		withIterator("-a", func(it *Iterator) {
			_, err := it.ReadInt32()
			require.IsType(t, InvalidDigitError{}, err)
		})
	})
	t.Run("negative overflow", func(t *testing.T) {
		s := fmt.Sprint(-math.MaxInt32 - 2)
		withIterator(s, func(it *Iterator) {
			_, err := it.ReadInt32()
			e, ok := err.(IntOverflowError)
			require.True(t, ok)
			require.Equal(t, "int32", e.typ)
			require.Equal(t, s, e.value)
		})
	})
	t.Run("negative max", func(t *testing.T) {
		s := fmt.Sprint(-math.MaxInt32 - 1)
		withIterator(s, func(it *Iterator) {
			i, err := it.ReadInt32()
			require.NoError(t, err)
			require.Equal(t, int32(-math.MaxInt32-1), i)
		})
	})
	t.Run("positive readUint32 error", func(t *testing.T) {
		withIterator("a", func(it *Iterator) {
			_, err := it.ReadInt32()
			require.IsType(t, InvalidDigitError{}, err)
		})
	})
	t.Run("positive overflow", func(t *testing.T) {
		s := fmt.Sprint(math.MaxInt32 + 1)
		withIterator(s, func(it *Iterator) {
			_, err := it.ReadInt32()
			e, ok := err.(IntOverflowError)
			require.True(t, ok)
			require.Equal(t, "int32", e.typ)
			require.Equal(t, s, e.value)
		})
	})
	t.Run("positive max", func(t *testing.T) {
		s := fmt.Sprint(math.MaxInt32)
		withIterator(s, func(it *Iterator) {
			i, err := it.ReadInt32()
			require.NoError(t, err)
			require.Equal(t, int32(math.MaxInt32), i)
		})
	})
}
