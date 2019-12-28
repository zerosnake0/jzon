package jzon

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Float_ReadFloat32(t *testing.T) {
	t.Run("leading space", func(t *testing.T) {
		withIterator(" 1", func(it *Iterator) {
			f, err := it.ReadFloat32()
			require.NoError(t, err)
			require.Equal(t, float32(1), f)
		})
	})
	t.Run("first byte error", func(t *testing.T) {
		withIterator(" ", func(it *Iterator) {
			_, err := it.ReadFloat32()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("negative error", func(t *testing.T) {
		withIterator("-", func(it *Iterator) {
			_, err := it.ReadFloat32()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("negative", func(t *testing.T) {
		withIterator("-0", func(it *Iterator) {
			f, err := it.ReadFloat32()
			require.NoError(t, err)
			require.Equal(t, float32(0), f)
		})
	})
	t.Run("negative 1", func(t *testing.T) {
		withIterator("-3.1415926535", func(it *Iterator) {
			f, err := it.ReadFloat32()
			require.NoError(t, err)
			require.InDelta(t, float32(-3.1415926535), f, 1e-10)
		})
	})
	t.Run("positive", func(t *testing.T) {
		withIterator("0", func(it *Iterator) {
			f, err := it.ReadFloat32()
			require.NoError(t, err)
			require.Equal(t, float32(0), f)
		})
	})
	t.Run("positive 1", func(t *testing.T) {
		withIterator("3.1415926535", func(it *Iterator) {
			f, err := it.ReadFloat32()
			require.NoError(t, err)
			require.InDelta(t, float32(3.1415926535), f, 1e-10)
		})
	})
	t.Run("invalid first char", func(t *testing.T) {
		withIterator("a", func(it *Iterator) {
			_, err := it.ReadFloat32()
			require.IsType(t, InvalidFloatError{}, err)
		})
	})
}

func TestIterator_Float_ReadFloat32_LeadingZero(t *testing.T) {
	t.Run("reader error", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			e := errors.New("test")
			it.Reset(&oneByteReader{
				b:   "0",
				err: e,
			})
			_, err := it.ReadFloat32()
			require.Equal(t, e, err)
		})
	})
	t.Run("fraction eof after dot", func(t *testing.T) {
		withIterator("0.", func(it *Iterator) {
			it.ResetBytes([]byte("0."))
			_, err := it.ReadFloat32()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("fraction invalid byte after dot", func(t *testing.T) {
		withIterator("0.a", func(it *Iterator) {
			_, err := it.ReadFloat32()
			require.IsType(t, InvalidFloatError{}, err)
		})
	})
	t.Run("fraction 1", func(t *testing.T) {
		withIterator("0.1a", func(it *Iterator) {
			f, err := it.ReadFloat32()
			require.NoError(t, err)
			require.InDelta(t, float32(0.1), f, 1e-10)
		})
	})
	t.Run("fraction 2", func(t *testing.T) {
		withIterator("0.1", func(it *Iterator) {
			f, err := it.ReadFloat32()
			require.NoError(t, err)
			require.InDelta(t, float32(0.1), f, 1e-10)
		})
	})
	t.Run("fraction reader error", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			e := errors.New("test")
			it.Reset(&oneByteReader{
				b:   "0.",
				err: e,
			})
			_, err := it.ReadFloat32()
			require.Equal(t, e, err)
		})
	})
	t.Run("fraction reader error 2", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			e := errors.New("test")
			it.Reset(&oneByteReader{
				b:   "0.1",
				err: e,
			})
			_, err := it.ReadFloat32()
			require.Equal(t, e, err)
		})
	})
	t.Run("exponent eof", func(t *testing.T) {
		withIterator("0.1e", func(it *Iterator) {
			_, err := it.ReadFloat32()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("exponent eof after sign", func(t *testing.T) {
		withIterator("0.1e+", func(it *Iterator) {
			_, err := it.ReadFloat32()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("exponent invalid byte after sign", func(t *testing.T) {
		withIterator("0.1e++", func(it *Iterator) {
			_, err := it.ReadFloat32()
			require.IsType(t, InvalidFloatError{}, err)
		})
	})
	t.Run("exponent 1", func(t *testing.T) {
		withIterator("0.1e+1", func(it *Iterator) {
			f, err := it.ReadFloat32()
			require.NoError(t, err)
			require.InDelta(t, float32(1), f, 1e-10)
		})
	})
	t.Run("exponent 2", func(t *testing.T) {
		withIterator("0.1e+1a", func(it *Iterator) {
			f, err := it.ReadFloat32()
			require.NoError(t, err)
			require.InDelta(t, float32(1), f, 1e-10)
		})
	})
	t.Run("exponent reader error", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			e := errors.New("test")
			it.Reset(&oneByteReader{
				b:   "0.1e+1",
				err: e,
			})
			_, err := it.ReadFloat32()
			require.Equal(t, e, err)
		})
	})
	t.Run("no fraction part", func(t *testing.T) {
		withIterator("0e+1", func(it *Iterator) {
			f, err := it.ReadFloat32()
			require.NoError(t, err)
			require.Equal(t, float32(0), f)
		})
	})
	t.Run("only zero", func(t *testing.T) {
		withIterator("0a", func(it *Iterator) {
			f, err := it.ReadFloat32()
			require.NoError(t, err)
			require.Equal(t, float32(0), f)
		})
	})
	t.Run("double zero", func(t *testing.T) {
		withIterator("00", func(it *Iterator) {
			f, err := it.ReadFloat32()
			require.NoError(t, err)
			require.Equal(t, float32(0), f)
		})
	})
}

func TestIterator_Float_ReadFloat32_NonLeadingZero(t *testing.T) {
	t.Run("exponent", func(t *testing.T) {
		withIterator("1e1", func(it *Iterator) {
			f, err := it.ReadFloat32()
			require.NoError(t, err)
			require.InDelta(t, float32(10), f, 1e-10)
		})
	})
	t.Run("no fraction no exponent 1", func(t *testing.T) {
		withIterator("10", func(it *Iterator) {
			f, err := it.ReadFloat32()
			require.NoError(t, err)
			require.InDelta(t, float32(10), f, 1e-10)
		})
	})
	t.Run("no fraction no exponent 2", func(t *testing.T) {
		withIterator("10a", func(it *Iterator) {
			f, err := it.ReadFloat32()
			require.NoError(t, err)
			require.InDelta(t, float32(10), f, 1e-10)
		})
	})
	t.Run("reader error", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			e := errors.New("test")
			it.Reset(&oneByteReader{
				b:   "10",
				err: e,
			})
			_, err := it.ReadFloat32()
			require.Equal(t, e, err)
		})
	})
}
