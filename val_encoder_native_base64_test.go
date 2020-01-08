package jzon

import (
	"testing"
)

// byte
type testMarshalByte byte

func (tb testMarshalByte) MarshalJSON() ([]byte, error) {
	return []byte{'"', '1', byte(tb), '"'}, nil
}

type testMarshalByte2 byte

func (tb *testMarshalByte2) MarshalJSON() ([]byte, error) {
	return []byte{'"', '2', byte(*tb), '"'}, nil
}

type testMarshalByte3 byte

func (tb testMarshalByte3) MarshalText() ([]byte, error) {
	return []byte{'"', '3', byte(tb), '"'}, nil
}

type testMarshalByte4 byte

func (tb *testMarshalByte4) MarshalText() ([]byte, error) {
	return []byte{'"', '4', byte(*tb), '"'}, nil
}

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
