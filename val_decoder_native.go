package jzon

import (
	"unsafe"
)

// bool decoder
type boolDecoder struct {
}

func (*boolDecoder) Decode(ptr unsafe.Pointer, it *Iterator) error {
	c, _, err := it.nextToken()
	if err != nil {
		return err
	}
	switch c {
	case 'n':
		it.head += 1
		return it.expectBytes("ull")
	case 't':
		it.head += 1
		if err := it.expectBytes("rue"); err != nil {
			return err
		}
		*(*bool)(ptr) = true
		return nil
	case 'f':
		it.head += 1
		if err := it.expectBytes("alse"); err != nil {
			return err
		}
		*(*bool)(ptr) = false
		return nil
	default:
		return UnexpectedByteError{got: c}
	}
}

// string decoder
type stringDecoder struct {
}

func (*stringDecoder) Decode(ptr unsafe.Pointer, it *Iterator) error {
	c, _, err := it.nextToken()
	if err != nil {
		return err
	}
	switch c {
	case '"':
		it.head += 1
		s, err := it.readString()
		if err != nil {
			return err
		}
		*(*string)(ptr) = s
		return nil
	case 'n':
		it.head += 1
		return it.expectBytes("ull")
	default:
		return UnexpectedByteError{got: c, exp: '"', exp2: 'n'}
	}
}
