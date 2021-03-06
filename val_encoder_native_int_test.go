package jzon

import (
	"math"
	"testing"
)

func TestValEncoder_Int(t *testing.T) {
	f := func(t *testing.T, i int) {
		checkEncodeValueWithStandard(t, i, nil)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxInt32)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, (*int)(nil), nil)
	})
	t.Run("pointer", func(t *testing.T) {
		i := int(math.MaxInt32)
		checkEncodeValueWithStandard(t, &i, nil)
	})
}

func TestValEncoder_Int8(t *testing.T) {
	f := func(t *testing.T, i int8) {
		checkEncodeValueWithStandard(t, i, nil)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxInt8)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, (*int8)(nil), nil)
	})
	t.Run("pointer", func(t *testing.T) {
		i := int8(math.MaxInt8)
		checkEncodeValueWithStandard(t, &i, nil)
	})
}

func TestValEncoder_Int16(t *testing.T) {
	f := func(t *testing.T, i int16) {
		checkEncodeValueWithStandard(t, i, nil)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxInt16)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, (*int16)(nil), nil)
	})
	t.Run("pointer", func(t *testing.T) {
		i := int16(math.MaxInt16)
		checkEncodeValueWithStandard(t, &i, nil)
	})
}

func TestValEncoder_Int32(t *testing.T) {
	f := func(t *testing.T, i int32) {
		checkEncodeValueWithStandard(t, i, nil)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxInt32)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, (*int32)(nil), nil)
	})
	t.Run("pointer", func(t *testing.T) {
		i := int32(math.MaxInt32)
		checkEncodeValueWithStandard(t, &i, nil)
	})
}

func TestValEncoder_Int64(t *testing.T) {
	f := func(t *testing.T, i int64) {
		checkEncodeValueWithStandard(t, i, nil)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxInt64)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, (*int64)(nil), nil)
	})
	t.Run("pointer", func(t *testing.T) {
		i := int64(math.MaxInt64)
		checkEncodeValueWithStandard(t, &i, nil)
	})
}

func TestValEncoder_Int_OmitEmpty(t *testing.T) {
	type st struct {
		I   int   `json:",omitempty"`
		I8  int8  `json:",omitempty"`
		I16 int16 `json:",omitempty"`
		I32 int32 `json:",omitempty"`
		I64 int64 `json:",omitempty"`
	}
	checkEncodeValueWithStandard(t, st{}, nil)
}
