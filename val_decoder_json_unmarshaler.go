package jzon

import (
	"encoding/json"
	"reflect"
	"unsafe"
)

var (
	jsonUnmarshalerType = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
)

type jsonUnmarshalerDecoder rtype

func (dec jsonUnmarshalerDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	obj := packEFace(rtype(dec), ptr)
	unmarshaler := obj.(json.Unmarshaler)
	b, err := it.SkipRaw()
	if err != nil {
		return err
	}
	return unmarshaler.UnmarshalJSON(b)
}
