package jzon

import (
	"math"
	"testing"
)

func TestValEncoder_Float32(t *testing.T) {
	f := func(t *testing.T, f float32) {
		checkEncodeValueWithStandard(t, DefaultEncoder, f, nil)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxFloat32)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, DefaultEncoder, (*float32)(nil), nil)
	})
	t.Run("pointer", func(t *testing.T) {
		f := float32(math.MaxFloat32)
		checkEncodeValueWithStandard(t, DefaultEncoder, &f, nil)
	})
}

func TestValEncoder_Float64(t *testing.T) {
	f := func(t *testing.T, f float64) {
		checkEncodeValueWithStandard(t, DefaultEncoder, f, nil)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxFloat64)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, DefaultEncoder, (*float64)(nil), nil)
	})
	t.Run("pointer", func(t *testing.T) {
		f := float64(math.MaxFloat64)
		checkEncodeValueWithStandard(t, DefaultEncoder, &f, nil)
	})
}
