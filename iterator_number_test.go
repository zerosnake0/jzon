package jzon

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Number_ReadNumber(t *testing.T) {
	t.Run("eof", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			_, err := it.ReadNumber()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("invalid type", func(t *testing.T) {
		withIterator(" +", func(it *Iterator) {
			_, err := it.ReadNumber()
			require.IsType(t, UnexpectedByteError{}, err)
		})
	})
	t.Run("negative eof", func(t *testing.T) {
		withIterator(" -", func(it *Iterator) {
			_, err := it.ReadNumber()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("invalid first byte after dash", func(t *testing.T) {
		withIterator(" --", func(it *Iterator) {
			_, err := it.ReadNumber()
			require.IsType(t, InvalidDigitError{}, err)
		})
	})
}

func TestIterator_Number_ReadNumber_LeadingZero(t *testing.T) {
	t.Run("one zero", func(t *testing.T) {
		withIterator(" 0", func(it *Iterator) {
			n, err := it.ReadNumber()
			require.NoError(t, err)
			require.Equal(t, "0", string(n))
		})
	})
	t.Run("reader error", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			e := errors.New("test")
			it.Reset(&stepByteReader{
				b:   " 0",
				err: e,
			})
			_, err := it.ReadNumber()
			require.Equal(t, e, err)
		})
	})
	t.Run("double zero", func(t *testing.T) {
		withIterator(" 00", func(it *Iterator) {
			n, err := it.ReadNumber()
			require.NoError(t, err)
			require.Equal(t, "0", string(n))
		})
	})
	t.Run("fraction eof", func(t *testing.T) {
		withIterator(" 0.", func(it *Iterator) {
			_, err := it.ReadNumber()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("invalid fraction", func(t *testing.T) {
		withIterator(" 0.a", func(it *Iterator) {
			_, err := it.ReadNumber()
			require.IsType(t, InvalidDigitError{}, err)
		})
	})
	t.Run("fraction end with eof", func(t *testing.T) {
		withIterator(" 0.12", func(it *Iterator) {
			n, err := it.ReadNumber()
			require.NoError(t, err)
			require.Equal(t, "0.12", string(n))
		})
	})
	t.Run("fraction end with other char", func(t *testing.T) {
		withIterator(" 0.12 ", func(it *Iterator) {
			n, err := it.ReadNumber()
			require.NoError(t, err)
			require.Equal(t, "0.12", string(n))
		})
	})
	t.Run("fraction error", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			e := errors.New("test")
			it.Reset(&stepByteReader{
				b:   " 0.12",
				err: e,
			})
			_, err := it.ReadNumber()
			require.Equal(t, e, err)
		})
	})
	t.Run("exponent eof", func(t *testing.T) {
		withIterator(" 0.1e", func(it *Iterator) {
			_, err := it.ReadNumber()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("exponent eof after sign", func(t *testing.T) {
		withIterator(" 0.1e+", func(it *Iterator) {
			_, err := it.ReadNumber()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("exponent invalid byte after sign", func(t *testing.T) {
		withIterator(" 0.1e++", func(it *Iterator) {
			_, err := it.ReadNumber()
			require.IsType(t, InvalidDigitError{}, err)
		})
	})
	t.Run("exponent end with eof", func(t *testing.T) {
		withIterator(" 0.1e+2", func(it *Iterator) {
			n, err := it.ReadNumber()
			require.NoError(t, err)
			require.Equal(t, "0.1e+2", string(n))
		})
	})
	t.Run("exponent end with another char", func(t *testing.T) {
		withIterator(" 0.1e+2 ", func(it *Iterator) {
			n, err := it.ReadNumber()
			require.NoError(t, err)
			require.Equal(t, "0.1e+2", string(n))
		})
	})
	t.Run("exponent end error", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			e := errors.New("test")
			it.Reset(&stepByteReader{
				b:   " 0.1e+2",
				err: e,
			})
			_, err := it.ReadNumber()
			require.Equal(t, e, err)
		})
	})
	t.Run("exponent only", func(t *testing.T) {
		withIterator(" 0e+1 ", func(it *Iterator) {
			n, err := it.ReadNumber()
			require.NoError(t, err)
			require.Equal(t, "0e+1", string(n))
		})
	})
}

func TestIterator_Number_ReadNumber_NonLeadingZero(t *testing.T) {
	t.Run("integer", func(t *testing.T) {
		withIterator(" 1", func(it *Iterator) {
			n, err := it.ReadNumber()
			require.NoError(t, err)
			require.Equal(t, "1", string(n))
		})
	})
	t.Run("reader error", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			e := errors.New("test")
			it.Reset(&stepByteReader{
				b:   " 1",
				err: e,
			})
			_, err := it.ReadNumber()
			require.Equal(t, e, err)
		})
	})
	t.Run("double digit", func(t *testing.T) {
		withIterator(" 12", func(it *Iterator) {
			n, err := it.ReadNumber()
			require.NoError(t, err)
			require.Equal(t, "12", string(n))
		})
	})
	t.Run("double digit end with other char", func(t *testing.T) {
		withIterator(" 12 ", func(it *Iterator) {
			n, err := it.ReadNumber()
			require.NoError(t, err)
			require.Equal(t, "12", string(n))
		})
	})
	t.Run("fraction eof", func(t *testing.T) {
		withIterator(" 1.", func(it *Iterator) {
			_, err := it.ReadNumber()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("invalid fraction", func(t *testing.T) {
		withIterator(" 1.+", func(it *Iterator) {
			_, err := it.ReadNumber()
			require.IsType(t, InvalidDigitError{}, err)
		})
	})
	t.Run("fraction end with eof", func(t *testing.T) {
		withIterator(" 1.23", func(it *Iterator) {
			n, err := it.ReadNumber()
			require.NoError(t, err)
			require.Equal(t, "1.23", string(n))
		})
	})
	t.Run("fraction end with other char", func(t *testing.T) {
		withIterator(" 1.23 ", func(it *Iterator) {
			n, err := it.ReadNumber()
			require.NoError(t, err)
			require.Equal(t, "1.23", string(n))
		})
	})
	t.Run("fraction error", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			e := errors.New("test")
			it.Reset(&stepByteReader{
				b:   " 1.23",
				err: e,
			})
			_, err := it.ReadNumber()
			require.Equal(t, e, err)
		})
	})
	t.Run("exponent eof", func(t *testing.T) {
		withIterator(" 1.2e", func(it *Iterator) {
			_, err := it.ReadNumber()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("exponent eof after sign", func(t *testing.T) {
		withIterator(" 1.2e+", func(it *Iterator) {
			_, err := it.ReadNumber()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("exponent invalid byte after sign", func(t *testing.T) {
		withIterator(" 1.2e++", func(it *Iterator) {
			_, err := it.ReadNumber()
			require.IsType(t, InvalidDigitError{}, err)
		})
	})
	t.Run("exponent end with eof", func(t *testing.T) {
		withIterator(" 1.2e+3", func(it *Iterator) {
			n, err := it.ReadNumber()
			require.NoError(t, err)
			require.Equal(t, "1.2e+3", string(n))
		})
	})
	t.Run("exponent end with another char", func(t *testing.T) {
		withIterator(" 1.2e+3 ", func(it *Iterator) {
			n, err := it.ReadNumber()
			require.NoError(t, err)
			require.Equal(t, "1.2e+3", string(n))
		})
	})
	t.Run("exponent end error", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			e := errors.New("test")
			it.Reset(&stepByteReader{
				b:   " 1.2e+3",
				err: e,
			})
			_, err := it.ReadNumber()
			require.Equal(t, e, err)
		})
	})
	t.Run("exponent only", func(t *testing.T) {
		withIterator(" 1e+2 ", func(it *Iterator) {
			n, err := it.ReadNumber()
			require.NoError(t, err)
			require.Equal(t, "1e+2", string(n))
		})
	})
}

func TestIterator_Number_ReadNumber_LargeNumber(t *testing.T) {
	s := "-" + strings.Repeat("123", 20) + "." +
		strings.Repeat("0456", 20) + "e+" +
		strings.Repeat("789", 20)
	withIterator(" "+s+" ", func(it *Iterator) {
		n, err := it.ReadNumber()
		require.NoError(t, err)
		require.Equal(t, s, string(n))
	})
}
