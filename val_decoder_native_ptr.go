package jzon

import (
	"reflect"
	"unsafe"
)

type pointerDecoder struct {
	elemDec   ValDecoder
	ptrRType  rtype
	elemRType rtype
}

func newPointerDecoder(typ reflect.Type) *pointerDecoder {
	return &pointerDecoder{
		ptrRType:  rtypeOfType(typ),
		elemRType: rtypeOfType(typ.Elem()),
	}
}

func (dec *pointerDecoder) Decode(ptr unsafe.Pointer, it *Iterator) error {
	c, _, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		if err = it.expectBytes("ull"); err != nil {
			return err
		}
		*(*unsafe.Pointer)(ptr) = nil
	} else {
		elemPtr := *((*unsafe.Pointer)(ptr))
		var tgtPtr unsafe.Pointer
		if elemPtr == nil {
			tgtPtr = unsafe_New(dec.elemRType)
		} else {
			tgtPtr = elemPtr
		}
		if err = dec.elemDec.Decode(tgtPtr, it); err != nil {
			return err
		}
		if elemPtr == nil {
			*(*unsafe.Pointer)(ptr) = tgtPtr
		}
	}
	return nil
}
