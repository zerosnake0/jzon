package jzon

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Skip_SkipObject(t *testing.T) {
	t.Run("eof", func(t *testing.T) {
		it := NewIterator()
		err := it.SkipObject()
		require.Equal(t, io.EOF, err)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" }"))
		err := it.SkipObject()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("eof after bracket", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" {"))
		err := it.SkipObject()
		require.Equal(t, io.EOF, err)
	})
	t.Run("empty", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" { } "))
		err := it.SkipObject()
		require.NoError(t, err)
		_, err = it.NextValueType()
		require.Equal(t, io.EOF, err)
	})
	t.Run("invalid field first byte", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { 1`))
		err := it.SkipObject()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("field error", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " `))
		err := it.SkipObject()
		require.Equal(t, io.EOF, err)
	})
	t.Run("value eof", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : `))
		err := it.SkipObject()
		require.Equal(t, io.EOF, err)
	})
	t.Run("invalid first value", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : +`))
		err := it.SkipObject()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("eof after first value", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : 1 `))
		err := it.SkipObject()
		require.Equal(t, io.EOF, err)
	})
	t.Run("only one pair", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : 1 } `))
		err := it.SkipObject()
		require.NoError(t, err)
	})
	t.Run("non nested second item dot error", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : true a `))
		err := it.SkipObject()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("non nested second item eof", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : false , `))
		err := it.SkipObject()
		require.Equal(t, io.EOF, err)
	})
	t.Run("non nested second item field error", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : false , a`))
		err := it.SkipObject()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("non nested second item", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : false , " " : true } `))
		err := it.SkipObject()
		require.NoError(t, err)
	})
	t.Run("nested eof", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : { `))
		err := it.SkipObject()
		require.Equal(t, io.EOF, err)
	})
	t.Run("nested empty value", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : { } } `))
		err := it.SkipObject()
		require.NoError(t, err)
	})
	t.Run("nested second item comma error", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : { } a `))
		err := it.SkipObject()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("nested second item no field", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : { } , `))
		err := it.SkipObject()
		require.Equal(t, io.EOF, err)
	})
	t.Run("nested second item bad field", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : { } , a`))
		err := it.SkipObject()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("nested second item field eof", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : { } , " " `))
		err := it.SkipObject()
		require.Equal(t, io.EOF, err)
	})
	t.Run("nested second item eof value", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : { } , " " : `))
		err := it.SkipObject()
		require.Equal(t, io.EOF, err)
	})
	t.Run("nested second item bad value", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : { } , " " : } `))
		err := it.SkipObject()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("nested second item", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : { } , " " : 0 } `))
		err := it.SkipObject()
		require.NoError(t, err)
		_, err = it.NextValueType()
		require.Equal(t, io.EOF, err)
	})
	t.Run("nested", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(nestedObject(100)))
		err := it.SkipObject()
		require.NoError(t, err)
		_, err = it.NextValueType()
		require.Equal(t, io.EOF, err)
	})
	t.Run("nested with array", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(nestedObjectWithArray(100)))
		err := it.SkipObject()
		require.NoError(t, err)
		_, err = it.NextValueType()
		require.Equal(t, io.EOF, err)
	})
}
