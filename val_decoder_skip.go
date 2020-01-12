package jzon

import (
	"unsafe"
)

type skipDecoder struct{}

func (*skipDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	it.head += 1
	if c != '{' && c != 'n' {
		return UnexpectedByteError{got: c, exp: '{', exp2: 'n'}
	}
	return skipFunctions[c](it, c)
}

type emptyObjectDecoder struct{}

func (*emptyObjectDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	it.head += 1
	switch c {
	case 'n':
		return it.expectBytes("ull")
	case '{':
		c, err := it.nextToken()
		if err != nil {
			return err
		}
		it.head += 1
		if c != '}' {
			return UnexpectedByteError{got: c, exp: '}'}
		}
		return nil
	default:
		return UnexpectedByteError{got: c, exp: '{', exp2: 'n'}
	}
}
