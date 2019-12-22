package jzon

import (
	"math"
	"testing"
)

func TestValEncoder_Float32(t *testing.T) {
	f := func(t *testing.T, f float32) {
		checkEncodeValueWithStandard(t, DefaultEncoder, f)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxFloat32)
	})
}

func TestValEncoder_Float64(t *testing.T) {
	f := func(t *testing.T, f float64) {
		checkEncodeValueWithStandard(t, DefaultEncoder, f)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxFloat64)
	})
}
