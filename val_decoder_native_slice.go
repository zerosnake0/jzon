package jzon

import (
	"encoding/base64"
	"reflect"
	"unsafe"
)

type sliceDecoderBuilder struct {
	decoder      *sliceDecoder
	elemPtrRType rtype
}

func newSliceDecoder(sliceType reflect.Type) *sliceDecoderBuilder {
	elem := sliceType.Elem()
	return &sliceDecoderBuilder{
		decoder: &sliceDecoder{
			elemKind:  elem.Kind(),
			elemRType: rtypeOfType(elem),
			elemSize:  elem.Size(),
		},
		elemPtrRType: rtypeOfType(reflect.PtrTo(elem)),
	}
}

func (builder *sliceDecoderBuilder) build(cache decoderCache) {
	builder.decoder.elemDec = cache[builder.elemPtrRType]
}

type sliceDecoder struct {
	elemKind  reflect.Kind
	elemRType rtype
	elemSize  uintptr
	elemDec   ValDecoder
}

func (dec *sliceDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		if err = it.expectBytes("ull"); err != nil {
			return err
		}
		sh := (*sliceHeader)(ptr)
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
		data, err := base64.StdEncoding.DecodeString(localByteToString(buf))
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
	c, err = it.nextToken()
	if err != nil {
		return err
	}
	newPtr := unsafeMakeSlice(dec.elemRType, 0, 0)
	if c == ']' {
		it.head += 1
	} else {
		for length := 1; ; length++ {
			newPtr = unsafeGrowSlice(dec.elemRType, newPtr, length)
			// must get the address every time
			childPtr := unsafeSliceChildPtr(newPtr, dec.elemSize, length-1)
			if err = dec.elemDec.Decode(childPtr, it, nil); err != nil {
				return err
			}
			c, err = it.nextToken()
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
	*(*sliceHeader)(ptr) = *(*sliceHeader)(newPtr)
	return nil
}
