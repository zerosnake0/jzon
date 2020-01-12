package jzon

import "unsafe"

type jsonNumberDecoder struct{}

func (*jsonNumberDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	var s string
	switch valueTypeMap[c] {
	case StringValue:
		it.head += 1
		s, err = it.readString()
		if err != nil {
			return err
		}
		*((*string)(ptr)) = s
		return nil
	case NullValue:
		// to be compatible with standard lib
		it.head += 1
		return it.expectBytes("ull")
	case NumberValue:
		it.head += 1
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
