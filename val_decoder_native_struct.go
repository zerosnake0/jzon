package jzon

import (
	"reflect"
	"strings"
	"unsafe"
)

type fieldInfo struct {
	offset  uintptr
	ptrType reflect.Type
	rtype   rtype
	decoder ValDecoder
}

type structDecoder struct {
	fields map[string]*fieldInfo
}

func (dec *Decoder) newStructDecoder(typ reflect.Type) *structDecoder {
	var key string
	var fields map[string]*fieldInfo
	for i := 0; i < typ.NumField(); i++ {
		stField := typ.Field(i)
		// field name cannot be empty (?)
		if stField.Name[0] < 'A' || stField.Name[0] > 'Z' {
			continue
		}
		tagV, ok := stField.Tag.Lookup(dec.tag)
		if ok {
			if tagV == "-" {
				continue
			}
			// TODO: complete this
			key = tagV
		} else {
			key = stField.Name
		}
		if !dec.caseSensitive {
			key = strings.ToLower(key)
		}
		if fields == nil {
			fields = map[string]*fieldInfo{}
		}
		fieldPtrType := reflect.PtrTo(stField.Type)
		fields[key] = &fieldInfo{
			offset:  stField.Offset,
			ptrType: fieldPtrType,
			rtype:   rtypeOfType(fieldPtrType),
		}
	}
	if len(fields) == 0 {
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
		_, field, err := it.readObjectFieldAsSlice(it.tmpBuffer[:0],
			it.decoder.caseSensitive)
		it.tmpBuffer = field
		if err != nil {
			return err
		}
		stField := dec.fields[localByteToString(field)]
		if stField != nil {
			fieldPtr := add(ptr, stField.offset, "struct field")
			if err = stField.decoder.Decode(fieldPtr, it); err != nil {
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
