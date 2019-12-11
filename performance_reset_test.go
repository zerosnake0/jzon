package jzon

import (
	"testing"
)

func resetBytes(it *Iterator, data []byte) {
	if it.reader != nil && it.buffer != nil {
		releaseByteSlice(it.buffer)
	}
	it.reader = nil
	it.buffer = data
	it.offset = 0
	it.head = 0
	it.tail = len(data)
}

func resetBytes2(it *Iterator, data []byte) {
	if it.reader != nil && it.buffer != nil {
		releaseByteSlice(it.buffer)
	}
	*it = Iterator{
		buffer: data,
		tail:   len(data),
	}
}

func Benchmark_Performance_Reset(b *testing.B) {
	data := make([]byte, 128)
	b.Run("impl", func(b *testing.B) {
		b.ReportAllocs()
		it := NewIterator()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			resetBytes(it, data)
		}
	})
	b.Run("alter", func(b *testing.B) {
		b.ReportAllocs()
		it := NewIterator()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			resetBytes2(it, data)
		}
	})
}
