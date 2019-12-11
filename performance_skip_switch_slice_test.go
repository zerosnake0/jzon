package jzon

import "testing"

func skip(it *Iterator, c byte, vt ValueType) error {
	switch c {
	case '"':
		return skipString(it, c)
	case 'n':
		return it.expectBytes("ull")
	case 't':
		return it.expectBytes("rue")
	case 'f':
		return it.expectBytes("alse")
	case '[':
		return skipArrayWithStack(it, c)
	case '{':
		return skipObjectWithStack(it, c)
	default:
		if vt != NumberValue {
			return UnexpectedByteError{got: c}
		}
		return skipNumber(it, c)
	}
}

func skip2(it *Iterator, c byte, _ ValueType) error {
	return skipFunctions[c](it, c)
}

func Benchmark_Performance_Skip_Switch(b *testing.B) {
	data := []byte(` "s", -123.0456e789, true, false, null, { } ]`)
	b.Run("impl", func(b *testing.B) {
		b.ReportAllocs()
		it := NewIterator()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			it.ResetBytes(data)
			skip2(it, '[', ArrayValue)
		}
	})
	b.Run("alter", func(b *testing.B) {
		b.ReportAllocs()
		it := NewIterator()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			it.ResetBytes(data)
			skip(it, '[', ArrayValue)
		}
	})
}
