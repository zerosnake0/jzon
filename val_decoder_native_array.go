package jzon

import (
	"reflect"
	"unsafe"
)

type arrayDecoderBuilder struct {
	decoder      *arrayDecoder
	elemPtrRType rtype
}

func newArrayDecoder(arrType reflect.Type) *arrayDecoderBuilder {
	elem := arrType.Elem()
	return &arrayDecoderBuilder{
		decoder: &arrayDecoder{
			rtype:    rtypeOfType(arrType),
			elemSize: elem.Size(),
			length:   arrType.Len(),
			// elemRType: rtypeOfType(elem),
		},
		elemPtrRType: rtypeOfType(reflect.PtrTo(elem)),
	}

}

type arrayDecoder struct {
	rtype    rtype
	elemSize uintptr
	length   int
	// elemRType rtype

	elemDec ValDecoder
}

func (dec *arrayDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	c, _, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		return it.expectBytes("ull")
	}
	if c != '[' {
		return UnexpectedByteError{got: c, exp: '[', exp2: 'n'}
	}
	it.head += 1
	c, _, err = it.nextToken()
	if err != nil {
		return err
	}
	count := 0
	childPtr := uintptr(ptr)
	if c == ']' {
		it.head += 1
	} else {
		for {
			if count < dec.length {
				elemPtr := unsafe.Pointer(childPtr)
				if err := dec.elemDec.Decode(elemPtr, it, nil); err != nil {
					return err
				}
				count += 1
				childPtr += dec.elemSize
			} else {
				if err := it.Skip(); err != nil {
					return err
				}
			}
			c, _, err = it.nextToken()
			if err != nil {
				return err
			}
			it.head += 1
			if c == ']' {
				break
			}
			if c != ',' {
				return UnexpectedByteError{got: c, exp: ']', exp2: ','}
			}
		}
	}
	if count < dec.length {
		// should be safe (?)
		typedmemclrpartial(dec.rtype, unsafe.Pointer(childPtr),
			uintptr(count)*dec.elemSize,
			uintptr(dec.length-count)*dec.elemSize)
	}
	return nil
}
