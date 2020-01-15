package jzon

import (
	"unsafe"
)

type structDecoderBuilder struct {
	decoder *structDecoder
	fields  structFields
}

func newStructDecoder(fields structFields) *structDecoderBuilder {
	return &structDecoderBuilder{
		decoder: &structDecoder{},
		fields:  fields,
	}
}

func (builder *structDecoderBuilder) build(cache decoderCache) {
	builder.decoder.fields.init(len(builder.fields))
	for i := range builder.fields {
		fi := &builder.fields[i]
		fiPtrRType := rtypeOfType(fi.ptrType)
		builder.decoder.fields.add(fi, cache[fiPtrRType])
	}
}

type decoderFieldInfo struct {
	offsets   []offset
	nameBytes []byte                 // []byte(name)
	equalFold func(s, t []byte) bool // bytes.EqualFold or equivalent
	quoted    bool
	decoder   ValDecoder
}

func (dfi *decoderFieldInfo) initFrom(f *field) {
	dfi.offsets = f.offsets
	dfi.nameBytes = f.nameBytes
	dfi.equalFold = f.equalFold
	dfi.quoted = f.quoted
}

type decoderFields struct {
	list           []decoderFieldInfo
	nameIndex      map[string]int
	nameIndexUpper map[string]int
}

func (df *decoderFields) init(size int) {
	df.list = make([]decoderFieldInfo, 0, size)
	df.nameIndex = make(map[string]int, size)
	df.nameIndexUpper = make(map[string]int, size)
}

func (df *decoderFields) add(f *field, dec ValDecoder) {
	df.nameIndex[f.name] = len(df.list)
	nameUpper := string(f.nameBytesUpper)
	if _, ok := df.nameIndexUpper[nameUpper]; !ok {
		df.nameIndexUpper[nameUpper] = len(df.list)
	}
	var dfi decoderFieldInfo
	dfi.initFrom(f)
	dfi.decoder = dec
	df.list = append(df.list, dfi)
}

func (df *decoderFields) find(key, buf []byte, caseSensitive bool) (*decoderFieldInfo, []byte) {
	if i, ok := df.nameIndex[localByteToString(key)]; ok {
		return &df.list[i], buf
	}
	if caseSensitive {
		return nil, buf
	}
	l := len(buf)
	// TODO: compare performance
	if true {
		// use the same buffer
		upper := toUpper(key, buf)
		i, ok := df.nameIndexUpper[localByteToString(upper[l:])]
		if ok {
			return &df.list[i], upper
		}
		return nil, upper
	} else {
		for i := range df.list {
			ff := &df.list[i]
			if ff.equalFold(ff.nameBytes, key) {
				return ff, buf
			}
		}
		return nil, buf
	}
}

type structDecoder struct {
	fields decoderFields
}

func (dec *structDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) (err error) {
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
		stField, fieldOut := dec.fields.find(field, it.tmpBuffer, it.decoder.caseSensitive)
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
