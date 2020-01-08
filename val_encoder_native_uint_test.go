package jzon

import (
	"math"
	"testing"
)

func TestValEncoder_Uint(t *testing.T) {
	f := func(t *testing.T, i uint) {
		checkEncodeValueWithStandard(t, i, nil)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxUint32)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, (*uint)(nil), nil)
	})
	t.Run("pointer", func(t *testing.T) {
		i := uint(math.MaxUint32)
		checkEncodeValueWithStandard(t, &i, nil)
	})
}

func TestValEncoder_Uint8(t *testing.T) {
	f := func(t *testing.T, i uint8) {
		checkEncodeValueWithStandard(t, i, nil)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxUint8)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, (*uint8)(nil), nil)
	})
	t.Run("pointer", func(t *testing.T) {
		i := uint8(math.MaxUint8)
		checkEncodeValueWithStandard(t, &i, nil)
	})
}

func TestValEncoder_Uint16(t *testing.T) {
	f := func(t *testing.T, i uint16) {
		checkEncodeValueWithStandard(t, i, nil)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxUint16)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, (*uint16)(nil), nil)
	})
	t.Run("pointer", func(t *testing.T) {
		i := uint16(math.MaxUint16)
		checkEncodeValueWithStandard(t, &i, nil)
	})
}

func TestValEncoder_Uint32(t *testing.T) {
	f := func(t *testing.T, i uint32) {
		checkEncodeValueWithStandard(t, i, nil)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxUint32)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, (*uint32)(nil), nil)
	})
	t.Run("pointer", func(t *testing.T) {
		i := uint32(math.MaxUint32)
		checkEncodeValueWithStandard(t, &i, nil)
	})
}

func TestValEncoder_Uint64(t *testing.T) {
	f := func(t *testing.T, i uint64) {
		checkEncodeValueWithStandard(t, i, nil)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxUint64)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, (*uint64)(nil), nil)
	})
	t.Run("pointer", func(t *testing.T) {
		i := uint64(math.MaxUint64)
		checkEncodeValueWithStandard(t, &i, nil)
	})
}

func TestValEncoder_Uint_OmitEmpty(t *testing.T) {
	type st struct {
		U    uint    `json:",omitempty"`
		U8   uint8   `json:",omitempty"`
		U16  uint16  `json:",omitempty"`
		U32  uint32  `json:",omitempty"`
		U64  uint64  `json:",omitempty"`
		Uptr uintptr `json:",omitempty"`
	}
	checkEncodeValueWithStandard(t, st{}, nil)
}
