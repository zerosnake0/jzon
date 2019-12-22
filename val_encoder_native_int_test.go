package jzon

import (
	"math"
	"testing"
)

func TestValEncoder_Int(t *testing.T) {
	f := func(t *testing.T, i int) {
		checkEncodeValueWithStandard(t, DefaultEncoder, i)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxInt8)
	})
}

func TestValEncoder_Int8(t *testing.T) {
	f := func(t *testing.T, i int8) {
		checkEncodeValueWithStandard(t, DefaultEncoder, i)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxInt8)
	})
}

func TestValEncoder_Int16(t *testing.T) {
	f := func(t *testing.T, i int16) {
		checkEncodeValueWithStandard(t, DefaultEncoder, i)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxInt16)
	})
}

func TestValEncoder_Int32(t *testing.T) {
	f := func(t *testing.T, i int32) {
		checkEncodeValueWithStandard(t, DefaultEncoder, i)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxInt32)
	})
}

func TestValEncoder_Int64(t *testing.T) {
	f := func(t *testing.T, i int64) {
		checkEncodeValueWithStandard(t, DefaultEncoder, i)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxInt64)
	})
}
