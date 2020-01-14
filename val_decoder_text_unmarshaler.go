package jzon

import (
	"encoding"
	"reflect"
	"unsafe"
)

var (
	textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
)

type textUnmarshalerDecoder rtype

func (dec textUnmarshalerDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	obj := packEFace(rtype(dec), ptr)
	unmarshaler := obj.(encoding.TextUnmarshaler)
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	switch c {
	case '"':
		it.head += 1
		b, err := it.readStringAsSlice()
		if err != nil {
			return err
		}
		return unmarshaler.UnmarshalText(b)
	case 'n':
		it.head += 1
		return it.expectBytes("ull")
	default:
		return UnexpectedByteError{got: c, exp: '"', exp2: 'n'}
	}
}
