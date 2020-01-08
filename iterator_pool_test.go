package jzon

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIteratorPool(t *testing.T) {
	must := require.New(t)

	pool := NewIteratorPool()

	f := func(cb func(it *Iterator)) {
		it := pool.BorrowIterator()
		must.Nil(it.reader)
		must.Nil(it.buffer)
		must.Equal(0, it.offset)
		must.Equal(0, it.head)
		must.Equal(0, it.tail)
		cb(it)
		pool.ReturnIterator(it)
		must.Nil(it.reader)
		must.Nil(it.buffer)
	}

	f(func(it *Iterator) {
		it.Reset(&oneByteReader{})
	})

	data := []byte("test")

	f(func(it *Iterator) {
		it.Reset(bytes.NewBuffer(data))
		must.Nil(it.reader)
		must.Equal(data, it.buffer)
		must.Equal(0, it.offset)
		must.Equal(0, it.head)
		must.Equal(len(data), it.tail)
	})

	f(func(it *Iterator) {
		it.ResetBytes(data)
		must.Nil(it.reader)
		must.Equal(data, it.buffer)
		must.Equal(0, it.offset)
		must.Equal(0, it.head)
		must.Equal(len(data), it.tail)
	})
}
