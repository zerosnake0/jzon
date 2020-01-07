package jzon

import (
	"math"
	"testing"
)

func TestValEncoder_Float32(t *testing.T) {
	f := func(t *testing.T, f float32) {
		checkEncodeValueWithStandard(t, f, nil)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxFloat32)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, (*float32)(nil), nil)
	})
	t.Run("pointer", func(t *testing.T) {
		f := float32(math.MaxFloat32)
		checkEncodeValueWithStandard(t, &f, nil)
	})
}

func TestValEncoder_Float64(t *testing.T) {
	f := func(t *testing.T, f float64) {
		checkEncodeValueWithStandard(t, f, nil)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxFloat64)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, (*float64)(nil), nil)
	})
	t.Run("pointer", func(t *testing.T) {
		f := float64(math.MaxFloat64)
		checkEncodeValueWithStandard(t, &f, nil)
	})
}

func TestValEncoder_Float32_OmitEmpty(t *testing.T) {
	type st struct {
		A float32 `json:",omitempty"`
	}
	t.Run("zero", func(t *testing.T) {
		checkEncodeValueWithStandard(t, st{}, nil)
	})
	t.Run("explicit zero", func(t *testing.T) {
		checkEncodeValueWithStandard(t, st{
			A: 0.0,
		}, nil)
	})
}

func TestValEncoder_Float64_OmitEmpty(t *testing.T) {
	type st struct {
		A float64 `json:",omitempty"`
	}
	t.Run("zero", func(t *testing.T) {
		checkEncodeValueWithStandard(t, st{}, nil)
	})
	t.Run("explicit zero", func(t *testing.T) {
		checkEncodeValueWithStandard(t, st{
			A: 0.0,
		}, nil)
	})
}
