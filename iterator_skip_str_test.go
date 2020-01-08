package jzon

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Skip_SkipString(t *testing.T) {
	t.Run("eof", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			err := it.SkipString()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("invalid first byte", func(t *testing.T) {
		withIterator(` a`, func(it *Iterator) {
			err := it.SkipString()
			require.IsType(t, UnexpectedByteError{}, err)
		})
	})
	t.Run("empty string", func(t *testing.T) {
		withIterator(` ""`, func(it *Iterator) {
			err := it.SkipString()
			require.NoError(t, err)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("eof on escape", func(t *testing.T) {
		withIterator(` "\`, func(it *Iterator) {
			err := it.SkipString()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("valid escape", func(t *testing.T) {
		withIterator(` "\t" `, func(it *Iterator) {
			err := it.SkipString()
			require.NoError(t, err)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("invalid escape", func(t *testing.T) {
		withIterator(` "\a`, func(it *Iterator) {
			err := it.SkipString()
			require.IsType(t, InvalidEscapeCharError{}, err)
		})
	})
	t.Run("invalid u4", func(t *testing.T) {
		withIterator(` "\uX`, func(it *Iterator) {
			err := it.SkipString()
			require.IsType(t, InvalidUnicodeCharError{}, err)
		})
	})
	t.Run("eof u4", func(t *testing.T) {
		withIterator(` "\u0`, func(it *Iterator) {
			err := it.SkipString()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("valid u4", func(t *testing.T) {
		withIterator(` "\u0000" `, func(it *Iterator) {
			err := it.SkipString()
			require.NoError(t, err)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("invalid string char", func(t *testing.T) {
		s := []byte{' ', '"', 0}
		withIterator(localByteToString(s), func(it *Iterator) {
			err := it.SkipString()
			require.IsType(t, InvalidStringCharError{}, err)
		})
	})
	t.Run("eof after first byte", func(t *testing.T) {
		withIterator(` " `, func(it *Iterator) {
			err := it.SkipString()
			require.Equal(t, io.EOF, err)
		})
	})
}
