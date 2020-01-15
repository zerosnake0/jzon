package jzon

import (
	"bytes"
	"reflect"
	"unsafe"
)

type oneFieldStructDecoderBuilder struct {
	decoder      *oneFieldStructDecoder
	fieldPtrType reflect.Type
}

func newOneFieldStructDecoder(field *field, caseSensitive bool) *oneFieldStructDecoderBuilder {
	var dfi decoderFieldInfo
	dfi.initFrom(field)
	if caseSensitive {
		dfi.equalFold = bytes.Equal
	}
	return &oneFieldStructDecoderBuilder{
		decoder: &oneFieldStructDecoder{
			decoderFieldInfo: dfi,
		},
		fieldPtrType: field.ptrType,
	}
}

func (builder *oneFieldStructDecoderBuilder) build(cache decoderCache) {
	builder.decoder.decoder = cache[rtypeOfType(builder.fieldPtrType)]
}

type oneFieldStructDecoder struct {
	decoderFieldInfo
}

func (dec *oneFieldStructDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) (err error) {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		err = it.expectBytes("ull")
		return
	}
	if c != '{' {
		return UnexpectedByteError{got: c, exp2: 'n', exp: '{'}
	}
	it.head += 1
	c, err = it.nextToken()
	if err != nil {
		return
	}
	if c == '}' {
		it.head += 1
		return
	}
	if c != '"' {
		return UnexpectedByteError{got: c, exp: '}', exp2: '"'}
	}
	it.head += 1
	for {
		field, err := it.readObjectFieldAsSlice()
		if err != nil {
			return err
		}
		if dec.equalFold(field, dec.nameBytes) {
			curPtr := add(ptr, dec.offsets[0].val, "struct field")
			for _, offset := range dec.offsets[1:] {
				subPtr := *(*unsafe.Pointer)(curPtr)
				if subPtr == nil {
					if offset.rtype == 0 { // the ptr field is not exported
						return NilEmbeddedPointerError
					}
					subPtr = unsafe_New(offset.rtype)
					*(*unsafe.Pointer)(curPtr) = subPtr
				}
				curPtr = add(subPtr, offset.val, "struct field")
			}
			opt := DecOpts{
				Quoted: dec.quoted,
			}
			if err = dec.decoder.Decode(curPtr, it, opt.noescape()); err != nil {
				return err
			}
		} else {
			if it.decoder.disallowUnknownFields {
				return UnknownFieldError(field)
			}
			if err = it.Skip(); err != nil {
				return err
			}
		}
		c, err = it.nextToken()
		if err != nil {
			return err
		}
		switch c {
		case '}':
			it.head += 1
			return nil
		case ',':
			it.head += 1
			c, err = it.nextToken()
			if err != nil {
				return err
			}
			if c != '"' {
				return UnexpectedByteError{got: c, exp: '"'}
			}
			it.head += 1
		default:
			return UnexpectedByteError{got: c, exp: '}', exp2: ','}
		}
	}
}
