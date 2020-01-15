package jzon

import (
	"reflect"
	"unsafe"
)

type pointerDecoderBuilder struct {
	decoder  *pointerDecoder
	ptrRType rtype
}

func newPointerDecoder(typ reflect.Type) *pointerDecoderBuilder {
	return &pointerDecoderBuilder{
		decoder: &pointerDecoder{
			elemRType: rtypeOfType(typ.Elem()),
		},
		ptrRType: rtypeOfType(typ),
	}
}

func (builder *pointerDecoderBuilder) build(cache decoderCache) {
	builder.decoder.elemDec = cache[builder.ptrRType]
}

type pointerDecoder struct {
	elemRType rtype

	elemDec ValDecoder
}

func (dec *pointerDecoder) Decode(ptr unsafe.Pointer, it *Iterator, opts *DecOpts) error {
	c, err := it.nextToken()
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
		if err = dec.elemDec.Decode(tgtPtr, it, opts); err != nil {
			return err
		}
		if elemPtr == nil {
			*(*unsafe.Pointer)(ptr) = tgtPtr
		}
	}
	return nil
}
