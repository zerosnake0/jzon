package jzon

import (
	"testing"
)

func TestValEncoder_Native_Struct_OmitEmpty(t *testing.T) {
	t.Run("mix", func(t *testing.T) {
		type st struct {
			S    string      `json:",omitempty"`
			I8   int8        `json:",omitempty"`
			I16  int16       `json:",omitempty"`
			I32  int32       `json:",omitempty"`
			I64  int64       `json:",omitempty"`
			U8   uint8       `json:",omitempty"`
			U16  uint16      `json:",omitempty"`
			U32  uint32      `json:",omitempty"`
			U64  uint64      `json:",omitempty"`
			Uptr uintptr     `json:",omitempty"`
			PI8  *int8       `json:",omitempty"`
			O    interface{} `json:",omitempty"`
		}
		checkEncodeValueWithStandard(t, st{}, nil)
		i8 := int8(1)
		checkEncodeValueWithStandard(t, st{
			PI8: &i8,
			O:   &i8,
		}, nil)
	})
	t.Run("struct", func(t *testing.T) {
		type inner struct {
			A int `json:",omitempty"`
		}
		type outer struct {
			B inner `json:",omitempty"`
		}
		checkEncodeValueWithStandard(t, outer{}, nil)
		checkEncodeValueWithStandard(t, outer{
			B: inner{A: 1},
		}, nil)
	})
	t.Run("struct 2", func(t *testing.T) {
		type inner struct {
			A int `json:",omitempty"`
		}
		type outer struct {
			B *inner `json:",omitempty"`
		}
		checkEncodeValueWithStandard(t, outer{}, nil)
		checkEncodeValueWithStandard(t, outer{
			B: &inner{},
		}, nil)
		checkEncodeValueWithStandard(t, outer{
			B: &inner{A: 1},
		}, nil)
	})
}
