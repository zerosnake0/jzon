package jzon

import (
	"encoding/base64"
	"io"
	"testing"
)

func TestValDecoder_Native_Base64(t *testing.T) {
	f := func(t *testing.T, data string, ex error, initValue string) {
		var p1 *[]byte
		var p2 *[]byte
		if initValue != "" {
			b1 := append([]byte(nil), initValue...)
			p1 = &b1
			b2 := append([]byte(nil), initValue...)
			p2 = &b2
		}
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
	}
	f2 := func(t *testing.T, data string, ex error) {
		f(t, data, ex, "dummy")
	}
	t.Run("nil pointer", func(t *testing.T) {
		f(t, "null", NilPointerReceiverError, "")
	})
	t.Run("eof", func(t *testing.T) {
		f2(t, "", io.EOF)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		f2(t, `true`, UnexpectedByteError{})
	})
	t.Run("invalid null", func(t *testing.T) {
		f2(t, "nul", io.EOF)
	})
	t.Run("null", func(t *testing.T) {
		f2(t, "null", nil)
	})
	t.Run("invalid string", func(t *testing.T) {
		f2(t, `"abc`, io.EOF)
	})
	t.Run("invalid base64", func(t *testing.T) {
		f2(t, `"abc"`, base64.CorruptInputError(0))
	})
	t.Run("valid base64", func(t *testing.T) {
		f2(t, `"`+base64.StdEncoding.EncodeToString(
			[]byte("abc"))+`"`, nil)
	})
	t.Run("invalid array", func(t *testing.T) {
		f2(t, `[`, io.EOF)
	})
	t.Run("empty array", func(t *testing.T) {
		f2(t, `[]`, nil)
	})
	t.Run("invalid uint8", func(t *testing.T) {
		f2(t, `[256]`, IntOverflowError{})
	})
	t.Run("invalid read more", func(t *testing.T) {
		f2(t, `[1`, io.EOF)
	})
}

type testByte byte

func (tb *testByte) UnmarshalJSON(data []byte) error {
	*tb = testByte(data[0] - '0' + 1)
	return nil
}

func TestValDecoder_Native_Base64_Override(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
	}
	t.Run("array", func(t *testing.T) {
		arr1 := [...]testByte{'d', 'u', 'm', 'm', 'y'}
		arr2 := [...]testByte{'d', 'u', 'm', 'm', 'y'}
		f(t, `"`+base64.StdEncoding.EncodeToString(
			[]byte("abc"))+`"`, UnexpectedByteError{}, &arr1, &arr2)
	})
	t.Run("invalid element type", func(t *testing.T) {
		arr1 := []int8{'d', 'u', 'm', 'm', 'y'}
		arr2 := []int8{'d', 'u', 'm', 'm', 'y'}
		f(t, `"`+base64.StdEncoding.EncodeToString(
			[]byte("abc"))+`"`, UnexpectedByteError{}, &arr1, &arr2)
	})
	t.Run("invalid base64", func(t *testing.T) {
		arr1 := []testByte{'d', 'u', 'm', 'm', 'y'}
		arr2 := []testByte{'d', 'u', 'm', 'm', 'y'}
		f(t, `"abc"`, base64.CorruptInputError(0), &arr1, &arr2)
	})
	t.Run("valid base64", func(t *testing.T) {
		arr1 := []testByte{'d', 'u', 'm', 'm', 'y'}
		arr2 := []testByte{'d', 'u', 'm', 'm', 'y'}
		f(t, `"`+base64.StdEncoding.EncodeToString(
			[]byte("abc"))+`"`, nil, &arr1, &arr2)
	})
	t.Run("slice", func(t *testing.T) {
		arr1 := []testByte{'d', 'u', 'm', 'm', 'y'}
		arr2 := []testByte{'d', 'u', 'm', 'm', 'y'}
		f(t, ` [ 1, 2, 3 ] `, nil, &arr1, &arr2)
	})
}
