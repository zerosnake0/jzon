package jzon

import (
	"bytes"
	"unsafe"
)

type smallStructDecoderBuilder struct {
	decoder *smallStructDecoder
	fields  []field
}

func newSmallStructDecoder(fields []field) *smallStructDecoderBuilder {
	dfiList := make([]decoderFieldInfo, len(fields))
	for i := range fields {
		dfiList[i].initFrom(&fields[i])
	}
	return &smallStructDecoderBuilder{
		decoder: &smallStructDecoder{
			fields: dfiList,
		},
		fields: fields,
	}
}

func (builder *smallStructDecoderBuilder) build(cache decoderCache) {
	for i := range builder.fields {
		builder.decoder.fields[i].decoder = cache[rtypeOfType(builder.fields[i].ptrType)]
	}
}

type smallStructDecoder struct {
	fields []decoderFieldInfo
}

func (dec *smallStructDecoder) find(field []byte, caseSensitive bool) *decoderFieldInfo {
	for i := range dec.fields {
		ff := &dec.fields[i]
		if bytes.Equal(ff.nameBytes, field) {
			return ff
		}
	}
	if caseSensitive {
		return nil
	}
	for i := range dec.fields {
		ff := &dec.fields[i]
		if ff.equalFold(ff.nameBytes, field) {
			return ff
		}
	}
	return nil
}

func (dec *smallStructDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) (err error) {
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
		stField := dec.find(field, it.decoder.caseSensitive)
		if stField != nil {
			curPtr := add(ptr, stField.offsets[0].val, "struct field")
			for _, offset := range stField.offsets[1:] {
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
				Quoted: stField.quoted,
			}
			if err = stField.decoder.Decode(curPtr, it, opt.noescape()); err != nil {
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
