package jzon

import (
	"testing"
)

func TestValEncoder_Native_Struct_OmitEmpty(t *testing.T) {
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
}
