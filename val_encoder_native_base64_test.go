package jzon

import (
	"testing"
)

type testMarshalByte byte

func (tb testMarshalByte) MarshalJSON() ([]byte, error) {
	return []byte{'"', byte(tb), '"'}, nil
}

type testMarshalByte2 byte

func (tb *testMarshalByte2) MarshalJSON() ([]byte, error) {
	return []byte{'"', byte(*tb), '"'}, nil
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
	t.Run("marshaler", func(t *testing.T) {
		checkEncodeValueWithStandard(t, []testMarshalByte("test"), nil)
	})
	t.Run("pointer marshaler", func(t *testing.T) {
		checkEncodeValueWithStandard(t, []testMarshalByte2("test"), nil)
	})
}
