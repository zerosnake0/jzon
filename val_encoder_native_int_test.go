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
		f(t, math.MaxInt32)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, DefaultEncoder, (*int)(nil))
	})
	t.Run("pointer", func(t *testing.T) {
		i := int(math.MaxInt32)
		checkEncodeValueWithStandard(t, DefaultEncoder, &i)
	})
}

func TestValEncoder_Int8(t *testing.T) {
	f := func(t *testing.T, i int8) {
		checkEncodeValueWithStandard(t, DefaultEncoder, i)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxInt8)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, DefaultEncoder, (*int8)(nil))
	})
	t.Run("pointer", func(t *testing.T) {
		i := int8(math.MaxInt8)
		checkEncodeValueWithStandard(t, DefaultEncoder, &i)
	})
}

func TestValEncoder_Int16(t *testing.T) {
	f := func(t *testing.T, i int16) {
		checkEncodeValueWithStandard(t, DefaultEncoder, i)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxInt16)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, DefaultEncoder, (*int16)(nil))
	})
	t.Run("pointer", func(t *testing.T) {
		i := int16(math.MaxInt16)
		checkEncodeValueWithStandard(t, DefaultEncoder, &i)
	})
}

func TestValEncoder_Int32(t *testing.T) {
	f := func(t *testing.T, i int32) {
		checkEncodeValueWithStandard(t, DefaultEncoder, i)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxInt32)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, DefaultEncoder, (*int32)(nil))
	})
	t.Run("pointer", func(t *testing.T) {
		i := int32(math.MaxInt32)
		checkEncodeValueWithStandard(t, DefaultEncoder, &i)
	})
}

func TestValEncoder_Int64(t *testing.T) {
	f := func(t *testing.T, i int64) {
		checkEncodeValueWithStandard(t, DefaultEncoder, i)
	}
	t.Run("test", func(t *testing.T) {
		f(t, math.MaxInt64)
	})
	t.Run("nil pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, DefaultEncoder, (*int64)(nil))
	})
	t.Run("pointer", func(t *testing.T) {
		i := int64(math.MaxInt64)
		checkEncodeValueWithStandard(t, DefaultEncoder, &i)
	})
}
