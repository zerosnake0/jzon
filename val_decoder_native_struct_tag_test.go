package jzon

import (
	"testing"
)

func TestValDecoder_Native_Struct_Tag(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		checkDecodeWithStandard(t, DefaultDecoder, data, ex, p1, p2)
	}
	t.Run("unicode (u)", func(t *testing.T) {
		type st struct {
			A int `json:"\u4e2d"`
		}
		f(t, `{"中":1}`, nil, &st{}, &st{})
	})
	t.Run("unicode (U)", func(t *testing.T) {
		type st struct {
			A int `json:"\U00004e2d"`
		}
		f(t, `{"中":1}`, nil, &st{}, &st{})
	})
	t.Run("unicode (x)", func(t *testing.T) {
		type st struct {
			A int `json:"\xe4\xb8\xad"`
		}
		f(t, `{"中":1}`, nil, &st{}, &st{})
	})
	t.Run("simple letter fold", func(t *testing.T) {
		type st struct {
			A int `json:"ABC"`
		}
		f(t, `{"abc":1}`, nil, &st{}, &st{})
	})
	t.Run("ascii equal fold", func(t *testing.T) {
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
}
