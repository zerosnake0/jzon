package jzon

import (
	"unsafe"
)

type skipDecoder struct {
}

func (*skipDecoder) Decode(ptr unsafe.Pointer, it *Iterator) error {
	c, _, err := it.nextToken()
	if err != nil {
		return err
	}
	it.head += 1
	if c != '{' && c != 'n' {
		return UnexpectedByteError{got: c, exp: '{', exp2: 'n'}
	}
	return skipFunctions[c](it, c)
}
