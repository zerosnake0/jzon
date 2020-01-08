package jzon

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Skip_SkipNumber(t *testing.T) {
	t.Run("eof", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			err := it.SkipNumber()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("invalid first byte", func(t *testing.T) {
		withIterator(" +", func(it *Iterator) {
			err := it.SkipNumber()
			require.IsType(t, UnexpectedByteError{}, err)
		})
	})
	t.Run("eof after negative sign", func(t *testing.T) {
		withIterator(" -", func(it *Iterator) {
			err := it.SkipNumber()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("invalid char after negative sign", func(t *testing.T) {
		withIterator(" - 1 ", func(it *Iterator) {
			err := it.SkipNumber()
			require.IsType(t, InvalidDigitError{}, err)
		})
	})
	t.Run("zero", func(t *testing.T) {
		withIterator(" -0", func(it *Iterator) {
			err := it.SkipNumber()
			require.NoError(t, err)
		})
	})
	t.Run("reader error after zero", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			e := errors.New("test")
			it.Reset(&oneByteReader{
				b:   " -0",
				err: e,
			})
			err := it.SkipNumber()
			require.Equal(t, e, err)
		})
	})
	t.Run("double zero", func(t *testing.T) {
		withIterator(" -00 ", func(it *Iterator) {
			err := it.SkipNumber()
			require.NoError(t, err)
		})
	})
	t.Run("zero with fraction", func(t *testing.T) {
		withIterator(" -0.1 ", func(it *Iterator) {
			err := it.SkipNumber()
			require.NoError(t, err)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("zero with exponent", func(t *testing.T) {
		withIterator(" -0e+1 ", func(it *Iterator) {
			err := it.SkipNumber()
			require.NoError(t, err)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("non zero", func(t *testing.T) {
		withIterator(" 1 ", func(it *Iterator) {
			err := it.SkipNumber()
			require.NoError(t, err)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("non zero with fraction", func(t *testing.T) {
		withIterator(" 1.1 ", func(it *Iterator) {
			err := it.SkipNumber()
			require.NoError(t, err)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("non zero with exponent", func(t *testing.T) {
		withIterator(" 1e-1 ", func(it *Iterator) {
			err := it.SkipNumber()
			require.NoError(t, err)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("non zero eof", func(t *testing.T) {
		withIterator(" 1", func(it *Iterator) {
			err := it.SkipNumber()
			require.NoError(t, err)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("non zero reader error", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			e := errors.New("test")
			it.Reset(&oneByteReader{
				b:   " 1",
				err: e,
			})
			err := it.SkipNumber()
			require.Equal(t, e, err)
		})
	})
	t.Run("fraction empty", func(t *testing.T) {
		withIterator(" 0.", func(it *Iterator) {
			err := it.SkipNumber()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("fraction invalid", func(t *testing.T) {
		withIterator(" 0.+", func(it *Iterator) {
			err := it.SkipNumber()
			require.IsType(t, InvalidDigitError{}, err)
		})
	})
	t.Run("fraction with exponent", func(t *testing.T) {
		withIterator(" -1.2e+3 ", func(it *Iterator) {
			err := it.SkipNumber()
			require.NoError(t, err)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("fraction eof", func(t *testing.T) {
		withIterator(" -1.2", func(it *Iterator) {
			err := it.SkipNumber()
			require.NoError(t, err)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("fraction reader error", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			e := errors.New("test")
			it.Reset(&oneByteReader{
				b:   " 1.2",
				err: e,
			})
			err := it.SkipNumber()
			require.Equal(t, e, err)
		})
	})
	t.Run("exponent empty", func(t *testing.T) {
		withIterator(" 0e", func(it *Iterator) {
			err := it.SkipNumber()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("exponent eof after sign", func(t *testing.T) {
		withIterator(" 0e+", func(it *Iterator) {
			err := it.SkipNumber()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("exponent invalid char after sign", func(t *testing.T) {
		withIterator(" 0e++", func(it *Iterator) {
			err := it.SkipNumber()
			require.IsType(t, InvalidDigitError{}, err)
		})
	})
	t.Run("exponent eof", func(t *testing.T) {
		withIterator(" 0e+1", func(it *Iterator) {
			err := it.SkipNumber()
			require.NoError(t, err)
		})
	})
	t.Run("exponent reader error", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			e := errors.New("test")
			it.Reset(&oneByteReader{
				b:   " 0e+1",
				err: e,
			})
			err := it.SkipNumber()
			require.Equal(t, e, err)
		})
	})
}

func TestIterator_Skip_SkipNumber_LargeNumber(t *testing.T) {
	s := "-" + strings.Repeat("123", 20) + "." +
		strings.Repeat("0456", 20) + "e+" +
		strings.Repeat("789", 20)
	withIterator(" "+s+" ", func(it *Iterator) {
		err := it.SkipNumber()
		require.NoError(t, err)
	})
}
