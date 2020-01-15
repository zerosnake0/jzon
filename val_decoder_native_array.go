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
			size:     arrType.Size(),
			elemSize: elem.Size(),
			length:   arrType.Len(),
			// elemRType: rtypeOfType(elem),
		},
		elemPtrRType: rtypeOfType(reflect.PtrTo(elem)),
	}
}

func (builder *arrayDecoderBuilder) build(cache decoderCache) {
	builder.decoder.elemDec = cache[builder.elemPtrRType]
}

type arrayDecoder struct {
	rtype    rtype
	size     uintptr
	elemSize uintptr
	length   int
	// elemRType rtype

	elemDec ValDecoder
}

func (dec *arrayDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	c, err := it.nextToken()
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
	c, err = it.nextToken()
	if err != nil {
		return err
	}
	count := 0
	var offset uintptr = 0
	if c == ']' {
		it.head += 1
	} else {
		for {
			if count < dec.length {
				// newer golang version seems to disallow direct
				// uintptr to unsafe.Pointer convert
				elemPtr := add(ptr, offset, "count < dec.length")
				if err := dec.elemDec.Decode(elemPtr, it, nil); err != nil {
					return err
				}
				count += 1
				offset += dec.elemSize
			} else {
				if err := it.Skip(); err != nil {
					return err
				}
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
	if count < dec.length {
		// should be safe (?)
		typedmemclrpartial(dec.rtype, add(ptr, offset, "count < dec.length"),
			offset, dec.size-offset)
	}
	return nil
}
