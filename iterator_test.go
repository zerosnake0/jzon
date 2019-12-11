package jzon

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIterator_New(t *testing.T) {
	must := require.New(t)

	it := NewIterator()
	must.Nil(it.reader)
	must.Nil(it.buffer)
	must.Equal(0, it.offset)
	must.Equal(0, it.head)
	must.Equal(0, it.tail)
	// must.Nil(it.Error)
}

func TestIterator_Reset_Nil(t *testing.T) {
	must := require.New(t)
	it := NewIterator()
	must.Nil(it.reader)
	must.Nil(it.buffer)
	must.Equal(0, it.offset)
	must.Equal(0, it.head)
	must.Equal(0, it.tail)

	it.Reset(nil)
	must.Nil(it.reader)
	must.Nil(it.buffer)
	must.Equal(0, it.offset)
	must.Equal(0, it.head)
	must.Equal(0, it.tail)
}

func TestIterator_Reset(t *testing.T) {
	must := require.New(t)

	// nil -> reader
	it := NewIterator()
	r := bytes.NewReader(nil)
	it.Reset(r)
	must.Equal(r, it.reader)
	must.NotEmpty(it.buffer)
	must.Equal(0, it.head)
	must.Equal(0, it.tail)

	// reader -> reader
	addr := &it.buffer[0]
	r2 := bytes.NewReader(nil)
	it.Reset(r2)
	must.Equal(r2, it.reader)
	must.True(addr == &it.buffer[0])
	must.Equal(0, it.head)
	must.Equal(0, it.tail)

	// reader -> byte
	b := []byte("abc")
	it.ResetBytes(b)
	must.Nil(it.reader)
	must.True(&b[0] == &it.buffer[0])
	must.Equal(0, it.head)
	must.Equal(len(b), it.tail)

	// nil -> byte
	it = NewIterator()
	b2 := []byte("abc")
	it.ResetBytes(b2)
	must.Nil(it.reader)
	must.True(&b2[0] == &it.buffer[0])
	must.Equal(0, it.head)
	must.Equal(len(b2), it.tail)

	// byte -> byte
	b3 := []byte("defg")
	it.ResetBytes(b3)
	must.Nil(it.reader)
	must.True(&b3[0] == &it.buffer[0])
	must.Equal(0, it.head)
	must.Equal(len(b3), it.tail)

	// byte -> reader
	r3 := bytes.NewReader(nil)
	it.Reset(r3)
	must.Equal(r3, it.reader)
	must.Equal(0, it.head)
	must.Equal(0, it.tail)
}

func TestIterator_NextValueType(t *testing.T) {
	must := require.New(t)
	it := NewIterator()
	for c, typ := range valueTypeMap {
		it.ResetBytes([]byte{byte(c)})
		next, err := it.NextValueType()
		if typ == WhiteSpaceValue {
			require.Equal(t, io.EOF, err)
		} else {
			require.NoError(t, err)
			must.Equal(typ, next)
		}
	}
}
