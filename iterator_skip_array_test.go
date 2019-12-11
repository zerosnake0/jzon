package jzon

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Skip_SkipArray(t *testing.T) {
	t.Run("eof", func(t *testing.T) {
		it := NewIterator()
		err := it.SkipArray()
		require.Equal(t, io.EOF, err)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" ]"))
		err := it.SkipArray()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("eof after bracket", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" ["))
		err := it.SkipArray()
		require.Equal(t, io.EOF, err)
	})
	t.Run("empty array", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" [ ] "))
		err := it.SkipArray()
		require.NoError(t, err)
		_, err = it.NextValueType()
		require.Equal(t, io.EOF, err)
	})
	t.Run("invalid byte after bracket", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" [ ,"))
		err := it.SkipArray()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("eof after first element", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" [ 1 "))
		err := it.SkipArray()
		require.Equal(t, io.EOF, err)
	})
	t.Run("one element", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" [ 1 ] "))
		err := it.SkipArray()
		require.NoError(t, err)
	})
	t.Run("invalid byte after element", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" [ 1 [ "))
		err := it.SkipArray()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("eof after dot", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" [ 1 , "))
		err := it.SkipArray()
		require.Equal(t, io.EOF, err)
	})
	t.Run("two elements", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" [ 1 , null ] "))
		err := it.SkipArray()
		require.NoError(t, err)
	})
	t.Run("nested error eof 1", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` [ [ `))
		err := it.SkipArray()
		require.Equal(t, io.EOF, err)
	})
	t.Run("nested error 1", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` [ [ 1 a`))
		err := it.SkipArray()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("nested error eof 2", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` [ [ 1 , `))
		err := it.SkipArray()
		require.Equal(t, io.EOF, err)
	})
	t.Run("nested error token", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` [ [ 1 , +`))
		err := it.SkipArray()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("nested 1", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(nestedArray1(100)))
		err := it.SkipArray()
		require.NoError(t, err)
		_, err = it.NextValueType()
		require.Equal(t, io.EOF, err)
	})
	t.Run("nested 2", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(nestedArray2(100)))
		err := it.SkipArray()
		require.NoError(t, err)
		_, err = it.NextValueType()
		require.Equal(t, io.EOF, err)
	})
	t.Run("nested with object", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(nestedArrayWithObject(100)))
		err := it.SkipArray()
		require.NoError(t, err)
		_, err = it.NextValueType()
		require.Equal(t, io.EOF, err)
	})
}
