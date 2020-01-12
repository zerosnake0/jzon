package jzon

import (
	"testing"
)

func TestValEncoder_Native_Struct_Tag(t *testing.T) {
	t.Run("quote", func(t *testing.T) {
		t.Run("invalid tag", func(t *testing.T) {
			checkEncodeValueWithStandard(t, &struct {
				A int `json:"\""`
			}{
				A: 1,
			}, nil)
		})
	})
	t.Run("unicode", func(t *testing.T) {
		t.Run("u", func(t *testing.T) {
			checkEncodeValueWithStandard(t, &struct {
				A int `json:"\u4e2d"`
			}{
				A: 1,
			}, nil)
		})
		t.Run("U", func(t *testing.T) {
			checkEncodeValueWithStandard(t, &struct {
				A int `json:"\U00004e2d"`
			}{
				A: 1,
			}, nil)
		})
		t.Run("x", func(t *testing.T) {
			checkEncodeValueWithStandard(t, &struct {
				A int `json:"\xe4\xb8\xad"`
			}{
				A: 1,
			}, nil)
		})
	})
}
