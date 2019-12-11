package jzon

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIteratorPool(t *testing.T) {
	must := require.New(t)

	pool := NewIteratorPool()

	// borrow 1
	it := pool.BorrowIterator()
	must.Nil(it.reader)
	must.Nil(it.buffer)
	must.Equal(0, it.offset)
	must.Equal(0, it.head)
	must.Equal(0, it.tail)

	// reset reader
	it.Reset(&oneByteReader{})

	// return 1
	pool.ReturnIterator(it)
	must.Nil(it.reader)
	must.Nil(it.buffer)

	// borrow 2
	it = pool.BorrowIterator()
	must.Nil(it.reader)
	must.Nil(it.buffer)
	must.Equal(0, it.offset)
	must.Equal(0, it.head)
	must.Equal(0, it.tail)

	// reset bytes buffer
	data := []byte("test")
	it.Reset(bytes.NewBuffer(data))
	must.Nil(it.reader)
	must.Equal(data, it.buffer)
	must.Equal(0, it.offset)
	must.Equal(0, it.head)
	must.Equal(len(data), it.tail)

	// return 2
	pool.ReturnIterator(it)
	must.Nil(it.reader)
	must.Nil(it.buffer)

	// reset bytes
	it.ResetBytes(data)
	must.Nil(it.reader)
	must.Equal(data, it.buffer)
	must.Equal(0, it.offset)
	must.Equal(0, it.head)
	must.Equal(len(data), it.tail)

	// return 3
	pool.ReturnIterator(it)
	must.Nil(it.reader)
	must.Nil(it.buffer)
}
