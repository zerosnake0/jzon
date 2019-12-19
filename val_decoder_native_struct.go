package jzon

import (
	"reflect"
	"unsafe"
)

type fieldInfo struct {
	offset  uintptr
	ptrType reflect.Type
	rtype   rtype
	decoder ValDecoder
}

type structDecoder struct {
	// fields map[string]*fieldInfo
	fields structFields
}

func (dec *Decoder) newStructDecoder(typ reflect.Type) *structDecoder {
	fields := describeStruct(typ, dec.tag, dec.onlyTaggedField)
	if fields.count() == 0 {
		return nil
	}
	return &structDecoder{
		fields: fields,
	}
}

func (dec *structDecoder) Decode(ptr unsafe.Pointer, it *Iterator) (err error) {
	c, _, err := it.nextToken()
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
	c, _, err = it.nextToken()
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
		_, field, err := it.readObjectFieldAsSlice(it.tmpBuffer[:0])
		it.tmpBuffer = field
		if err != nil {
			return err
		}
		stField := dec.fields.find(field, it.decoder.caseSensitive)
		if stField != nil {
			curPtr := ptr
			for _, offset := range stField.offsets {
				curPtr = add(curPtr, offset.offset, "struct field")
				if offset.typ == nil {
					break
				}
				curPtr = *(*unsafe.Pointer)(curPtr)
				if curPtr == nil {
					return NilEmbeddedPointerError
				}
			}
			if err = stField.decoder.Decode(curPtr, it); err != nil {
				return err
			}
		} else {
			if err = it.Skip(); err != nil {
				return err
			}
		}
		c, _, err = it.nextToken()
		if err != nil {
			return err
		}
		switch c {
		case '}':
			it.head += 1
			return nil
		case ',':
			it.head += 1
			c, _, err = it.nextToken()
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
