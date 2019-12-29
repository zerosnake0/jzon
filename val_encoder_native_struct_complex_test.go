package jzon

import (
	"testing"
)

func TestValEncoder_Native_Struct_Complex_OneField(t *testing.T) {
	t.Run("field indirect", func(t *testing.T) {
		type st struct {
			A int
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, st{
			A: 1,
		}, nil)
	})
	t.Run("field direct", func(t *testing.T) {
		type st struct {
			A *int
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, st{
			A: nil,
		}, nil)
	})
	t.Run("direct (array)", func(t *testing.T) {
		type st struct {
			A [1]*int
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, st{
			A: [1]*int{nil},
		}, nil)
		i := 1
		checkEncodeValueWithStandard(t, DefaultEncoder, st{
			A: [1]*int{&i},
		}, nil)
	})
	t.Run("direct map", func(t *testing.T) {
		type st struct {
			A map[int]int
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, st{
			A: nil,
		}, nil)
		checkEncodeValueWithStandard(t, DefaultEncoder, st{
			A: map[int]int{1: 2},
		}, nil)
	})
}

func TestValEncoder_Native_Struct_Complex_MultipleField(t *testing.T) {
	t.Run("nil pointer", func(t *testing.T) {
		type st struct {
			A int
			B int
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, (*st)(nil), nil)
	})
	t.Run("pointer", func(t *testing.T) {
		type st struct {
			A int
			B int
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, &st{
			A: 1, B: 2,
		}, nil)
	})
	t.Run("non pointer", func(t *testing.T) {
		type st struct {
			A int
			B int
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, st{
			A: 1, B: 2,
		}, nil)
	})
	t.Run("nested", func(t *testing.T) {
		type inner struct {
			A int
			B int
		}
		type outer struct {
			*inner
			C int
		}
		type outer2 struct {
			inner
			C int
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, outer{
			inner: nil,
			C:     1,
		}, nil)
		checkEncodeValueWithStandard(t, DefaultEncoder, outer{
			inner: &inner{
				A: 1,
				B: 2,
			},
			C: 3,
		}, nil)
		checkEncodeValueWithStandard(t, DefaultEncoder, outer2{}, nil)
	})
	t.Run("pointer field", func(t *testing.T) {
		type st struct {
			A *int
			B *int
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, st{}, nil)
	})
}

func TestValEncoder_Native_Struct_Complex_Nested(t *testing.T) {
	t.Run("nested 1", func(t *testing.T) {
		type inner struct {
			A int
		}
		type outer struct {
			inner
			A int
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, outer{
			inner: inner{A: 1},
			A:     2,
		}, nil)
	})
	t.Run("nested 2", func(t *testing.T) {
		type inner struct {
			A int
		}
		type aliasInner inner
		type inner2inner struct {
			A int
			B int
		}
		type inner2 struct {
			inner2inner
		}
		type outer struct {
			inner
			aliasInner
			inner2
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, outer{
			inner: inner{
				A: 1,
			},
			aliasInner: aliasInner{
				A: 2,
			},
			inner2: inner2{
				inner2inner: inner2inner{
					A: 3,
					B: 4,
				},
			},
		}, nil)
	})
}
