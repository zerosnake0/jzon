package jzon

import (
	"unsafe"
)

type emptyObjectDecoder struct{}

func (*emptyObjectDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	it.head++
	switch c {
	case 'n':
		return it.expectBytes("ull")
	case '{':
		if it.disallowUnknownFields {
			c, err := it.nextToken()
			if err != nil {
				return err
			}
			it.head++
			if c != '}' {
				return UnexpectedByteError{got: c, exp: '}'}
			}
			return nil
		}
		return skipFunctions[c](it, c)
	default:
		return UnexpectedByteError{got: c, exp: '{', exp2: 'n'}
	}
}
