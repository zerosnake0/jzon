package jzon

import (
	"testing"
)

func TestValEncoder_Native_Struct_Tag(t *testing.T) {
	t.Run("quote (invalid tag)", func(t *testing.T) {
		checkEncodeValueWithStandard(t, &struct {
			A int `json:"\""`
		}{
			A: 1,
		}, nil)
	})
	t.Run("unicode (u)", func(t *testing.T) {
		checkEncodeValueWithStandard(t, &struct {
			A int `json:"\u4e2d"`
		}{
			A: 1,
		}, nil)
	})
	t.Run("unicode (U)", func(t *testing.T) {
		checkEncodeValueWithStandard(t, &struct {
			A int `json:"\U00004e2d"`
		}{
			A: 1,
		}, nil)
	})
	t.Run("unicode (x)", func(t *testing.T) {
		checkEncodeValueWithStandard(t, &struct {
			A int `json:"\xe4\xb8xad"`
		}{
			A: 1,
		}, nil)
	})
}
