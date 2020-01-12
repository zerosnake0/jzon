package jzon

import (
	"testing"
)

func TestValDecoder_Native_Struct_Tag(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		checkDecodeWithStandard(t, DefaultDecoder, data, ex, p1, p2)
	}
	t.Run("quote", func(t *testing.T) {
		t.Run("invalid tag", func(t *testing.T) {
			type st struct {
				A int `json:"\""`
			}
			f(t, `{"\"":1}`, nil, &st{}, &st{})
		})
	})
	t.Run("unicode", func(t *testing.T) {
		t.Run("u", func(t *testing.T) {
			type st struct {
				A int `json:"\u4e2d"`
			}
			f(t, `{"中":1}`, nil, &st{}, &st{})
		})
		t.Run("U", func(t *testing.T) {
			type st struct {
				A int `json:"\U00004e2d"`
			}
			f(t, `{"中":1}`, nil, &st{}, &st{})
		})
		t.Run("x", func(t *testing.T) {
			type st struct {
				A int `json:"\xe4\xb8\xad"`
			}
			f(t, `{"中":1}`, nil, &st{}, &st{})
		})
		t.Run("中", func(t *testing.T) {
			type st struct {
				A int `json:"中"`
			}
			f(t, `{"中":1}`, nil, &st{}, &st{})
		})
	})
	t.Run("fold", func(t *testing.T) {
		t.Run("simple letter", func(t *testing.T) {
			type st struct {
				A int `json:"ABC"`
			}
			f(t, `{"abc":1}`, nil, &st{}, &st{})
		})
		t.Run("ascii equal", func(t *testing.T) {
			type st struct {
				A int `json:"a0"`
			}
			f(t, `{"A0":1}`, nil, &st{}, &st{})
		})
		t.Run("equal fold right", func(t *testing.T) {
			type st struct {
				A int `json:"b"`
				B int `json:"ask"`
			}
			f(t, `{"AſK":1}`, nil, &st{}, &st{})
		})
	})
}
