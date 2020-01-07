package jzon

import (
	"testing"
)

func TestValEncoder_Direct_OmitEmpty(t *testing.T) {
	t.Run("array", func(t *testing.T) {
		type st struct {
			A [1]*int `json:",omitempty"`
		}
		checkEncodeValueWithStandard(t, &st{}, nil)
	})
}
