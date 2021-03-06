package jzon

import (
	"unsafe"
)

type jsonNumberDecoder struct{}

func (*jsonNumberDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	var s string
	switch valueTypeMap[c] {
	case StringValue:
		it.head++
		c, err = it.nextByte()
		if err != nil {
			return err
		}
		if valueTypeMap[c] != NumberValue {
			return InvalidDigitError{c}
		}
		it.head++
		s, err = it.readNumberAsString(c)
		if err != nil {
			return err
		}
		c, err = it.nextByte()
		if err != nil {
			return err
		}
		if c != '"' {
			return UnexpectedByteError{exp: '"', got: c}
		}
		it.head++
		*((*string)(ptr)) = s
		return nil
	case NullValue:
		// to be compatible with standard lib
		it.head++
		return it.expectBytes("ull")
	case NumberValue:
		it.head++
		s, err = it.readNumberAsString(c)
		if err != nil {
			return err
		}
		*((*string)(ptr)) = s
		return nil
	default:
		return UnexpectedByteError{got: c}
	}
}
