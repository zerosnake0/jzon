package jzon

import (
	"encoding/base64"
	"reflect"
	"unsafe"
)

type sliceDecoder struct {
	rtype        rtype
	elemKind     reflect.Kind
	elemDec      ValDecoder
	elemRType    rtype
	elemPtrRType rtype
	elemSize     uintptr
}

func newSliceDecoder(sliceType reflect.Type) *sliceDecoder {
	elem := sliceType.Elem()
	return &sliceDecoder{
		rtype:        rtypeOfType(sliceType),
		elemKind:     elem.Kind(),
		elemRType:    rtypeOfType(elem),
		elemPtrRType: rtypeOfType(reflect.PtrTo(elem)),
		elemSize:     elem.Size(),
	}
}

func (dec *sliceDecoder) Decode(ptr unsafe.Pointer, it *Iterator) error {
	c, _, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		if err = it.expectBytes("ull"); err != nil {
			return err
		}
		sh := (*reflect.SliceHeader)(ptr)
		sh.Len = 0
		sh.Cap = 0
		sh.Data = 0
		return nil
	}
	if c == '"' {
		if dec.elemKind != reflect.Uint8 {
			return UnexpectedByteError{got: c, exp: '[', exp2: 'n'}
		}
		it.head += 1
		// TODO: improve by checking only base64 characters
		begin := it.head
		oldCapture := it.capture
		it.capture = true
		if err := skipString(it, c); err != nil {
			return err
		}
		it.capture = oldCapture
		buf := it.buffer[begin : it.head-1]
		data, err := base64.StdEncoding.DecodeString(*(*string)(unsafe.Pointer(&buf)))
		if err != nil {
			return err
		}
		*((*[]byte)(ptr)) = data
		return nil
	}
	if c != '[' {
		return UnexpectedByteError{got: c, exp: '[', exp2: 'n'}
	}
	it.head += 1
	c, _, err = it.nextToken()
	if err != nil {
		return err
	}
	newPtr := unsafeMakeSlice(dec.rtype, 0, 0)
	if c == ']' {
		it.head += 1
	} else {
		for length := 1; ; length++ {
			newPtr = unsafeGrowSlice(dec.rtype, dec.elemRType, newPtr, length)
			// must get the address every time
			childPtr := unsafeSliceChildPtr(newPtr, dec.elemSize, length-1)
			if err = dec.elemDec.Decode(childPtr, it); err != nil {
				return err
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
	*(*reflect.SliceHeader)(ptr) = *(*reflect.SliceHeader)(newPtr)
	return nil
}
