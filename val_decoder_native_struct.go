package jzon

import (
	"reflect"
	"unsafe"
)

type structDecoderBuilder struct {
	decoder *structDecoder
	fields  structFields
}

type decoderFieldInfo struct {
	offsets []offset
	// nameBytes []byte                 // []byte(name)
	// equalFold func(s, t []byte) bool // bytes.EqualFold or equivalent
	quoted  bool
	decoder ValDecoder
}

type decoderFields struct {
	list           []decoderFieldInfo
	nameIndex      map[string]int
	nameIndexUpper map[string]int
}

func (df *decoderFields) init(size int) {
	df.list = make([]decoderFieldInfo, 0, size)
	df.nameIndex = make(map[string]int, size)
}

func (df *decoderFields) add(f *field, dec ValDecoder) {
	df.nameIndex[f.name] = len(df.list)
	nameUpper := string(f.nameBytesUpper)
	if _, ok := df.nameIndex[nameUpper]; !ok {
		df.nameIndex[nameUpper] = len(df.list)
	}
	df.list = append(df.list, decoderFieldInfo{
		offsets: f.offsets,
		// nameBytes: f.nameBytes,
		// equalFold: f.equalFold,
		quoted:  f.quoted,
		decoder: dec,
	})
}

func (df *decoderFields) find(key []byte, caseSensitive bool) (*decoderFieldInfo, []byte) {
	if i, ok := df.nameIndex[localByteToString(key)]; ok {
		return &df.list[i], key
	}
	if caseSensitive {
		return nil, key
	}
	l := len(key)
	// use the same buffer
	upper := toUpper(key, key)
	i, ok := df.nameIndexUpper[localByteToString(upper[l:])]
	if ok {
		return &df.list[i], upper
	}
	return nil, upper
	// // TODO: performance of this?
	// for i := range df.list {
	// 	ff := &df.list[i]
	// 	if ff.equalFold(ff.nameBytes, key) {
	// 		return ff
	// 	}
	// }
	// return nil
}

type structDecoder struct {
	fields decoderFields
}

func (dec *Decoder) newStructDecoder(typ reflect.Type) *structDecoderBuilder {
	fields := describeStruct(typ, dec.tag, dec.onlyTaggedField)
	if fields.count() == 0 {
		return nil
	}
	return &structDecoderBuilder{
		decoder: &structDecoder{},
		fields:  fields,
	}
}

func (dec *structDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) (err error) {
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
		stField, fieldOut := dec.fields.find(field, it.decoder.caseSensitive)
		it.tmpBuffer = fieldOut
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
			if err = stField.decoder.Decode(curPtr, it, &opt); err != nil {
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
