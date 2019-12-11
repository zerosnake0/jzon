package jzon

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Float_ReadFloat64(t *testing.T) {
	t.Run("leading space", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" 1"))
		f, err := it.ReadFloat64()
		require.NoError(t, err)
		require.Equal(t, float64(1), f)
	})
	t.Run("first byte error", func(t *testing.T) {
		it := NewIterator()
		_, err := it.ReadFloat64()
		require.Equal(t, io.EOF, err)
	})
	t.Run("negative error", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte{'-'})
		_, err := it.ReadFloat64()
		require.Equal(t, io.EOF, err)
	})
	t.Run("negative", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("-0"))
		f, err := it.ReadFloat64()
		require.NoError(t, err)
		require.Equal(t, float64(0), f)
	})
	t.Run("negative 1", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("-3.1415926535"))
		f, err := it.ReadFloat64()
		require.NoError(t, err)
		require.InDelta(t, float64(-3.1415926535), f, 1e-10)
	})
	t.Run("positive", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("0"))
		f, err := it.ReadFloat64()
		require.NoError(t, err)
		require.Equal(t, float64(0), f)
	})
	t.Run("positive 1", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("3.1415926535"))
		f, err := it.ReadFloat64()
		require.NoError(t, err)
		require.InDelta(t, float64(3.1415926535), f, 1e-10)
	})
	t.Run("invalid first char", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte{'a'})
		_, err := it.ReadFloat64()
		require.IsType(t, InvalidFloatError{}, err)
	})
}

func TestIterator_Float_ReadFloat64_LeadingZero(t *testing.T) {
	t.Run("reader error", func(t *testing.T) {
		it := NewIterator()
		e := errors.New("test")
		it.Reset(&oneByteReader{
			b:   []byte{'0'},
			err: e,
		})
		_, err := it.ReadFloat64()
		require.Equal(t, e, err)
	})
	t.Run("fraction eof after dot", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("0."))
		_, err := it.ReadFloat64()
		require.Equal(t, io.EOF, err)
	})
	t.Run("fraction invalid byte after dot", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("0.a"))
		_, err := it.ReadFloat64()
		require.IsType(t, InvalidFloatError{}, err)
	})
	t.Run("fraction 1", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("0.1a"))
		f, err := it.ReadFloat64()
		require.NoError(t, err)
		require.InDelta(t, float64(0.1), f, 1e-10)
	})
	t.Run("fraction 2", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("0.1"))
		f, err := it.ReadFloat64()
		require.NoError(t, err)
		require.InDelta(t, float64(0.1), f, 1e-10)
	})
	t.Run("fraction reader error", func(t *testing.T) {
		it := NewIterator()
		e := errors.New("test")
		it.Reset(&oneByteReader{
			b:   []byte("0."),
			err: e,
		})
		_, err := it.ReadFloat64()
		require.Equal(t, e, err)
	})
	t.Run("fraction reader error 2", func(t *testing.T) {
		it := NewIterator()
		e := errors.New("test")
		it.Reset(&oneByteReader{
			b:   []byte("0.1"),
			err: e,
		})
		_, err := it.ReadFloat64()
		require.Equal(t, e, err)
	})
	t.Run("exponent eof", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("0.1e"))
		_, err := it.ReadFloat64()
		require.Equal(t, io.EOF, err)
	})
	t.Run("exponent eof after sign", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("0.1e+"))
		_, err := it.ReadFloat64()
		require.Equal(t, io.EOF, err)
	})
	t.Run("exponent invalid byte after sign", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("0.1e++"))
		_, err := it.ReadFloat64()
		require.IsType(t, InvalidFloatError{}, err)
	})
	t.Run("exponent 1", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("0.1e+1"))
		f, err := it.ReadFloat64()
		require.NoError(t, err)
		require.InDelta(t, float64(1), f, 1e-10)
	})
	t.Run("exponent 2", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("0.1e+1a"))
		f, err := it.ReadFloat64()
		require.NoError(t, err)
		require.InDelta(t, float64(1), f, 1e-10)
	})
	t.Run("exponent reader error", func(t *testing.T) {
		it := NewIterator()
		e := errors.New("test")
		it.Reset(&oneByteReader{
			b:   []byte("0.1e+1"),
			err: e,
		})
		_, err := it.ReadFloat64()
		require.Equal(t, e, err)
	})
	t.Run("no fraction part", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("0e+1"))
		f, err := it.ReadFloat64()
		require.NoError(t, err)
		require.Equal(t, float64(0), f)
	})
	t.Run("only zero", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("0a"))
		f, err := it.ReadFloat64()
		require.NoError(t, err)
		require.Equal(t, float64(0), f)
	})
	t.Run("double zero", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("00"))
		f, err := it.ReadFloat64()
		require.NoError(t, err)
		require.Equal(t, float64(0), f)
	})
}

func TestIterator_Float_ReadFloat64_NonLeadingZero(t *testing.T) {
	t.Run("exponent", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("1e1"))
		f, err := it.ReadFloat64()
		require.NoError(t, err)
		require.InDelta(t, float64(10), f, 1e-10)
	})
	t.Run("no fraction no exponent 1", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("10"))
		f, err := it.ReadFloat64()
		require.NoError(t, err)
		require.InDelta(t, float64(10), f, 1e-10)
	})
	t.Run("no fraction no exponent 2", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte("10a"))
		f, err := it.ReadFloat64()
		require.NoError(t, err)
		require.InDelta(t, float64(10), f, 1e-10)
	})
	t.Run("reader error", func(t *testing.T) {
		it := NewIterator()
		e := errors.New("test")
		it.Reset(&oneByteReader{
			b:   []byte("10"),
			err: e,
		})
		_, err := it.ReadFloat64()
		require.Equal(t, e, err)
	})
}
