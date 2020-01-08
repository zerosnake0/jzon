package jzon

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

type testStringKind string

type testStringKindDecoder struct {
}

func (*testStringKindDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	s, err := it.ReadString()
	if err != nil {
		return err
	}
	*((*string)(ptr)) = "`" + s + "`"
	return nil
}

func TestValDecoder_Native_Kind_String(t *testing.T) {
	decoder := NewDecoder(&DecoderOption{
		ValDecoders: map[reflect.Type]ValDecoder{
			reflect.TypeOf(string("")): (*testStringKindDecoder)(nil),
		},
	})
	data := []byte(`"abc"`)
	var s testStringKind = "dummy"
	err := decoder.Unmarshal(data, &s)
	require.NoError(t, err)
	require.Equal(t, testStringKind("`abc`"), s)
}
