package jzon

import (
	"testing"
)

func resetTestImpl(it *Iterator) {
	it.reader = nil
	it.buffer = nil
	it.tail = 0

	it.capture = false
	it.offset = 0
	it.head = 0
	it.lastEfaceOffset = 0
	it.lastEfacePtr = 0
	it.Context = nil
}

// The current implementation is better in both:
// - maintainability
// - performance
func Benchmark_Performance_Reset(b *testing.B) {
	b.Run("impl", func(b *testing.B) {
		b.ReportAllocs()
		it := NewIterator()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			it.reset()
		}
	})
	b.Run("alter", func(b *testing.B) {
		b.ReportAllocs()
		it := NewIterator()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			resetTestImpl(it)
		}
	})
}
