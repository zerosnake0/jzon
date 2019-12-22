package jzon

import (
	"math"
	"testing"
)

func TestValEncoder_Uint(t *testing.T) {
	f := func(t *testing.T, i uint) {
		checkEncodeValueWithStandard(t, DefaultEncoder, i)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxUint8)
	})
}

func TestValEncoder_Uint8(t *testing.T) {
	f := func(t *testing.T, i uint8) {
		checkEncodeValueWithStandard(t, DefaultEncoder, i)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxUint8)
	})
}

func TestValEncoder_Uint16(t *testing.T) {
	f := func(t *testing.T, i uint16) {
		checkEncodeValueWithStandard(t, DefaultEncoder, i)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxUint16)
	})
}

func TestValEncoder_Uint32(t *testing.T) {
	f := func(t *testing.T, i uint32) {
		checkEncodeValueWithStandard(t, DefaultEncoder, i)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxUint32)
	})
}

func TestValEncoder_Uint64(t *testing.T) {
	f := func(t *testing.T, i uint64) {
		checkEncodeValueWithStandard(t, DefaultEncoder, i)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxUint64)
	})
}
