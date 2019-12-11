package jzon

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Object_ReadObjectBegin(t *testing.T) {
	t.Run("eof", func(t *testing.T) {
		it := NewIterator()
		_, _, err := it.ReadObjectBegin()
		require.Equal(t, io.EOF, err)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte{'"'})
		_, _, err := it.ReadObjectBegin()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("eof after first byte", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte{'{'})
		_, _, err := it.ReadObjectBegin()
		require.Equal(t, io.EOF, err)
	})
	t.Run("empty object", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" { } "))
		more, _, err := it.ReadObjectBegin()
		require.NoError(t, err)
		require.False(t, more)
	})
	t.Run("invalid", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { "" `))
		_, _, err := it.ReadObjectBegin()
		require.Equal(t, io.EOF, err)
	})
	t.Run("invalid colon", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { "" a `))
		_, _, err := it.ReadObjectBegin()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("valid", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { "" : `))
		more, field, err := it.ReadObjectBegin()
		require.NoError(t, err)
		require.True(t, more)
		require.Equal(t, "", field)
	})
	t.Run("invalid second token", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { 1 `))
		_, _, err := it.ReadObjectBegin()
		require.IsType(t, UnexpectedByteError{}, err)
	})
}

func TestIterator_Object_ReadObjectMore(t *testing.T) {
	init := func(t *testing.T, s string) *Iterator {
		it := NewIterator()
		buf := append([]byte(` { "k" : 2 `), s...)
		it.ResetBytes(buf)
		more, field, err := it.ReadObjectBegin()
		require.NoError(t, err)
		require.True(t, more)
		require.Equal(t, "k", field)

		i, err := it.ReadInt()
		require.NoError(t, err)
		require.Equal(t, 2, i)
		return it
	}
	t.Run("eof", func(t *testing.T) {
		it := init(t, "")
		_, _, err := it.ReadObjectMore()
		require.Equal(t, io.EOF, err)
	})
	t.Run("valid ending", func(t *testing.T) {
		it := init(t, "}")
		more, _, err := it.ReadObjectMore()
		require.NoError(t, err)
		require.False(t, more)
	})
	t.Run("eof after comma", func(t *testing.T) {
		it := init(t, `, `)
		_, _, err := it.ReadObjectMore()
		require.Equal(t, io.EOF, err)
	})
	t.Run("invalid byte after comma", func(t *testing.T) {
		it := init(t, `, a`)
		_, _, err := it.ReadObjectMore()
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("valid", func(t *testing.T) {
		it := init(t, `, "k2" : `)
		more, field, err := it.ReadObjectMore()
		require.NoError(t, err)
		require.True(t, more)
		require.Equal(t, "k2", field)
	})
	t.Run("invalid byte", func(t *testing.T) {
		it := init(t, ` a`)
		_, _, err := it.ReadObjectMore()
		require.IsType(t, UnexpectedByteError{}, err)
	})
}

func TestIterator_Object_ReadObject_Example(t *testing.T) {
	must := require.New(t)

	it := NewIterator()
	it.ResetBytes([]byte(` { "key" : "value" , "key2" : 1 } `))
	more, field, err := it.ReadObjectBegin()
	must.NoError(err)
	must.True(more)
	must.Equal("key", field)

	s, err := it.ReadString()
	must.NoError(err)
	must.Equal("value", s)

	more, field, err = it.ReadObjectMore()
	must.NoError(err)
	must.True(more)
	must.Equal("key2", field)

	i, err := it.ReadInt()
	must.NoError(err)
	must.Equal(1, i)

	more, _, err = it.ReadObjectMore()
	must.NoError(err)
	must.False(more)

	_, err = it.NextValueType()
	must.Equal(io.EOF, err)
}

func TestIterator_Object_ReadObjectCB(t *testing.T) {
	t.Run("eof", func(t *testing.T) {
		it := NewIterator()
		err := it.ReadObjectCB(nil)
		require.Equal(t, io.EOF, err)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" a"))
		err := it.ReadObjectCB(nil)
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("eof after first byte", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" { "))
		err := it.ReadObjectCB(nil)
		require.Equal(t, io.EOF, err)
	})
	t.Run("empty object", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(" { } "))
		err := it.ReadObjectCB(nil)
		require.NoError(t, err)
	})
	t.Run("invalid field", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { a`))
		err := it.ReadObjectCB(nil)
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("invalid field 2", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " `))
		err := it.ReadObjectCB(nil)
		require.Equal(t, io.EOF, err)
	})
	t.Run("error during callback", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : `))
		e := errors.New("test")
		err := it.ReadObjectCB(func(it *Iterator, field string) error {
			return e
		})
		require.Equal(t, e, err)
	})
	t.Run("eof after first item", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : 1`))
		err := it.ReadObjectCB(func(it *Iterator, field string) error {
			i, err := it.ReadInt()
			require.NoError(t, err)
			require.Equal(t, 1, i)
			return nil
		})
		require.Equal(t, io.EOF, err)
	})
	t.Run("end after first item", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : 1 } `))
		err := it.ReadObjectCB(func(it *Iterator, field string) error {
			i, err := it.ReadInt()
			require.NoError(t, err)
			require.Equal(t, 1, i)
			return nil
		})
		require.NoError(t, err)

		_, err = it.NextValueType()
		require.Equal(t, io.EOF, err)
	})
	t.Run("eof after comma", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : 1 , `))
		err := it.ReadObjectCB(func(it *Iterator, field string) error {
			i, err := it.ReadInt()
			require.NoError(t, err)
			require.Equal(t, 1, i)
			return nil
		})
		require.Equal(t, io.EOF, err)
	})
	t.Run("invalid byte after comma", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : 1 , a `))
		err := it.ReadObjectCB(func(it *Iterator, field string) error {
			i, err := it.ReadInt()
			require.NoError(t, err)
			require.Equal(t, 1, i)
			return nil
		})
		require.IsType(t, UnexpectedByteError{}, err)
	})
	t.Run("eof after second field", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : 1 , " `))
		err := it.ReadObjectCB(func(it *Iterator, field string) error {
			i, err := it.ReadInt()
			require.NoError(t, err)
			require.Equal(t, 1, i)
			return nil
		})
		require.Equal(t, io.EOF, err)
	})
	t.Run("invalid comma", func(t *testing.T) {
		it := NewIterator()
		it.ResetBytes([]byte(` { " " : 1 a " `))
		err := it.ReadObjectCB(func(it *Iterator, field string) error {
			i, err := it.ReadInt()
			require.NoError(t, err)
			require.Equal(t, 1, i)
			return nil
		})
		require.IsType(t, UnexpectedByteError{}, err)
	})
}

func TestIterator_Object_ReadObjectCB_Example(t *testing.T) {
	must := require.New(t)

	it := NewIterator()
	it.ResetBytes([]byte(` { "key" : "value" , "key2" : "value2" } `))

	m := map[string]string{}

	err := it.ReadObjectCB(func(it *Iterator, field string) (err error) {
		value, err := it.ReadString()
		if err == nil {
			m[field] = value
		}
		return
	})
	must.NoError(err)
	must.Len(m, 2)
	must.Equal("value", m["key"])
	must.Equal("value2", m["key2"])

	_, err = it.NextValueType()
	must.Equal(io.EOF, err)
}

func TestIterator_Object_skipObjectField(t *testing.T) {
	must := require.New(t)

	it := NewIterator()
	it.ResetBytes([]byte(` key" : `))
	more, err := it.skipObjectField()
	must.NoError(err)
	must.True(more)
}
