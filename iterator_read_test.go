package jzon

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_Read_Read(t *testing.T) {
	t.Run("eof", func(t *testing.T) {
		withIterator("", func(it *Iterator) {
			_, err := it.Read()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("string", func(t *testing.T) {
		withIterator(`"test"`, func(it *Iterator) {
			o, err := it.Read()
			require.NoError(t, err)
			require.Equal(t, "test", o)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("null", func(t *testing.T) {
		withIterator(`null`, func(it *Iterator) {
			o, err := it.Read()
			require.NoError(t, err)
			require.Equal(t, nil, o)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("true", func(t *testing.T) {
		withIterator(`true`, func(it *Iterator) {
			o, err := it.Read()
			require.NoError(t, err)
			require.Equal(t, true, o)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("false", func(t *testing.T) {
		withIterator(`false`, func(it *Iterator) {
			o, err := it.Read()
			require.NoError(t, err)
			require.Equal(t, false, o)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("number", func(t *testing.T) {
		withIterator(`-123.456e7`, func(it *Iterator) {
			o, err := it.Read()
			require.NoError(t, err)
			require.Equal(t, -123.456e7, o)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("invalid", func(t *testing.T) {
		withIterator(`+`, func(it *Iterator) {
			_, err := it.Read()
			require.IsType(t, UnexpectedByteError{}, err)
		})
	})
}

func TestIterator_Read_ReadArray(t *testing.T) {
	t.Run("eof after bracket", func(t *testing.T) {
		withIterator(`[`, func(it *Iterator) {
			_, err := it.Read()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("empty array", func(t *testing.T) {
		withIterator(`[]`, func(it *Iterator) {
			o, err := it.Read()
			require.NoError(t, err)
			require.Equal(t, []interface{}{}, o)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("bad element", func(t *testing.T) {
		withIterator(` [ + `, func(it *Iterator) {
			_, err := it.Read()
			require.IsType(t, UnexpectedByteError{}, err)
		})
	})
	t.Run("eof after first element", func(t *testing.T) {
		withIterator(` [ 1 `, func(it *Iterator) {
			_, err := it.Read()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("not nested", func(t *testing.T) {
		withIterator(` [ null ] `, func(it *Iterator) {
			o, err := it.Read()
			require.NoError(t, err)
			require.Equal(t, []interface{}{nil}, o)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("bad char after element", func(t *testing.T) {
		withIterator(` [ 1 + `, func(it *Iterator) {
			_, err := it.Read()
			require.IsType(t, UnexpectedByteError{}, err)
		})
	})
	t.Run("eof after comma", func(t *testing.T) {
		withIterator(` [ 1 , `, func(it *Iterator) {
			_, err := it.Read()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("mix elements", func(t *testing.T) {
		withIterator(` [ null , 1 , "a" , true , false ] `, func(it *Iterator) {
			o, err := it.Read()
			require.NoError(t, err)
			require.Equal(t, []interface{}{nil, float64(1), "a", true, false}, o)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("nested eof", func(t *testing.T) {
		withIterator(` [ [ `, func(it *Iterator) {
			_, err := it.Read()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("nested bad comma", func(t *testing.T) {
		withIterator(` [ [ ] [ `, func(it *Iterator) {
			_, err := it.Read()
			require.IsType(t, UnexpectedByteError{}, err)
		})
	})
	t.Run("nested eof after comma", func(t *testing.T) {
		withIterator(` [ [ ] , `, func(it *Iterator) {
			_, err := it.Read()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("nested bad item", func(t *testing.T) {
		withIterator(` [ [ ] , + `, func(it *Iterator) {
			_, err := it.Read()
			require.IsType(t, UnexpectedByteError{}, err)
		})
	})
	t.Run("nested 1", func(t *testing.T) {
		withIterator(` [ [ ] , null ] `, func(it *Iterator) {
			o, err := it.Read()
			require.NoError(t, err)
			t.Log(o)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("nested 2", func(t *testing.T) {
		s := nestedArray1(10)
		withIterator(s, func(it *Iterator) {
			o, err := it.Read()
			require.NoError(t, err)
			t.Log(o)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("nested 3", func(t *testing.T) {
		s := nestedArray2(10)
		withIterator(s, func(it *Iterator) {
			o, err := it.Read()
			require.NoError(t, err)
			t.Log(o)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("nested with object", func(t *testing.T) {
		s := nestedArrayWithObject(10)
		withIterator(s, func(it *Iterator) {
			o, err := it.Read()
			require.NoError(t, err)
			t.Log(o)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
}

func TestIterator_Read_ReadObject(t *testing.T) {
	t.Run("eof after bracket", func(t *testing.T) {
		withIterator(`{`, func(it *Iterator) {
			_, err := it.Read()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("empty", func(t *testing.T) {
		withIterator(` { } `, func(it *Iterator) {
			o, err := it.Read()
			require.NoError(t, err)
			require.Equal(t, map[string]interface{}{}, o)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("invalid char after bracket", func(t *testing.T) {
		withIterator(`{ ,`, func(it *Iterator) {
			_, err := it.Read()
			require.IsType(t, UnexpectedByteError{}, err)
		})
	})
	t.Run("field eof", func(t *testing.T) {
		withIterator(`{ " `, func(it *Iterator) {
			_, err := it.Read()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("eof after colon", func(t *testing.T) {
		withIterator(`{ " " : `, func(it *Iterator) {
			_, err := it.Read()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("bad field value", func(t *testing.T) {
		withIterator(`{ " " : + `, func(it *Iterator) {
			_, err := it.Read()
			require.IsType(t, UnexpectedByteError{}, err)
		})
	})
	t.Run("eof after field value", func(t *testing.T) {
		withIterator(`{ " " : 1 `, func(it *Iterator) {
			_, err := it.Read()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("one field", func(t *testing.T) {
		withIterator(` { "k" : "v" } `, func(it *Iterator) {
			o, err := it.Read()
			require.NoError(t, err)
			require.Equal(t, map[string]interface{}{"k": "v"}, o)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("bad char after field value", func(t *testing.T) {
		withIterator(`{ " " : 1 { `, func(it *Iterator) {
			_, err := it.Read()
			require.IsType(t, UnexpectedByteError{}, err)
		})
	})
	t.Run("eof after comma", func(t *testing.T) {
		withIterator(`{ " " : 1 , `, func(it *Iterator) {
			_, err := it.Read()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("two fields", func(t *testing.T) {
		withIterator(` { "k1" : "v1", "k2" : "v2" } `, func(it *Iterator) {
			o, err := it.Read()
			require.NoError(t, err)
			require.Equal(t, map[string]interface{}{"k1": "v1", "k2": "v2"}, o)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("nested eof", func(t *testing.T) {
		withIterator(` { "a" : { `, func(it *Iterator) {
			_, err := it.Read()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("nested bad comma", func(t *testing.T) {
		withIterator(` { "a" : { } + `, func(it *Iterator) {
			_, err := it.Read()
			require.IsType(t, UnexpectedByteError{}, err)
		})
	})
	t.Run("nested eof after comma", func(t *testing.T) {
		withIterator(` { "a" : { } , `, func(it *Iterator) {
			_, err := it.Read()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("nested bad char after comma", func(t *testing.T) {
		withIterator(` { "a" : { } , { `, func(it *Iterator) {
			_, err := it.Read()
			require.IsType(t, UnexpectedByteError{}, err)
		})
	})
	t.Run("nested bad field", func(t *testing.T) {
		withIterator(` { "a" : { } , " `, func(it *Iterator) {
			_, err := it.Read()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("nested eof after field", func(t *testing.T) {
		withIterator(` { "a" : { } , "b" : `, func(it *Iterator) {
			_, err := it.Read()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("nested bad value", func(t *testing.T) {
		withIterator(` { "a" : { } , "b" : } `, func(it *Iterator) {
			_, err := it.Read()
			require.IsType(t, UnexpectedByteError{}, err)
		})
	})
	t.Run("nested 1", func(t *testing.T) {
		withIterator(` { "a" : { } , "b" : null } `, func(it *Iterator) {
			o, err := it.Read()
			require.NoError(t, err)
			t.Log(o)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("nested 2", func(t *testing.T) {
		s := nestedObject(10)
		withIterator(s, func(it *Iterator) {
			o, err := it.Read()
			require.NoError(t, err)
			t.Log(o)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
	t.Run("nested with array", func(t *testing.T) {
		s := nestedObjectWithArray(10)
		withIterator(s, func(it *Iterator) {
			o, err := it.Read()
			require.NoError(t, err)
			t.Log(o)
			_, err = it.NextValueType()
			require.Equal(t, io.EOF, err)
		})
	})
}

func TestIterator_Read_UseNumber(t *testing.T) {
	dec := NewDecoder(&DecoderOption{
		UseNumber: true,
	})
	it := dec.NewIterator()
	defer dec.ReturnIterator(it)
	it.ResetBytes([]byte(" 123 "))
	o, err := it.Read()
	require.NoError(t, err)
	require.Equal(t, json.Number("123"), o)
	t.Log(o)
	_, err = it.NextValueType()
	require.Equal(t, io.EOF, err)
}
