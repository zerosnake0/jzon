package jzon

import (
	"testing"
)

func TestValEncoder_Native_Base64(t *testing.T) {
	t.Run("nil byte", func(t *testing.T) {
		checkEncodeValueWithStandard(t, []byte(nil), nil)
	})
	t.Run("empty byte", func(t *testing.T) {
		checkEncodeValueWithStandard(t, []byte{}, nil)
	})
	t.Run("type alias", func(t *testing.T) {
		type B byte
		checkEncodeValueWithStandard(t, []B("test"), nil)
	})
	t.Run("json marshaler", func(t *testing.T) {
		t.Run("value", func(t *testing.T) {
			checkEncodeValueWithStandard(t, []testMarshalByte("test"), nil)
		})
		t.Run("pointer", func(t *testing.T) {
			checkEncodeValueWithStandard(t, []testMarshalByte2("test"), nil)
		})
	})
	t.Run("text marshaler", func(t *testing.T) {
		t.Run("value", func(t *testing.T) {
			checkEncodeValueWithStandard(t, []testMarshalByte3("test"), nil)
		})
		t.Run("pointer", func(t *testing.T) {
			checkEncodeValueWithStandard(t, []testMarshalByte4("test"), nil)
		})
	})
}

func TestValEncoder_Native_Base64_OmitEmpty(t *testing.T) {
	type st struct {
		A []byte `json:",omitempty"`
	}
	t.Run("nil", func(t *testing.T) {
		checkEncodeValueWithStandard(t, st{}, nil)
	})
	t.Run("empty", func(t *testing.T) {
		checkEncodeValueWithStandard(t, st{
			A: []byte{},
		}, nil)
	})
	t.Run("non empty", func(t *testing.T) {
		checkEncodeValueWithStandard(t, st{
			A: []byte("test"),
		}, nil)
	})
}
