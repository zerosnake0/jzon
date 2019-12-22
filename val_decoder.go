package jzon

import (
	"encoding/json"
	"reflect"
	"strconv"
	"unsafe"
)

var (
	globalValDecoders = map[rtype]ValDecoder{}
	kindMap           = [numKinds]rtype{}
	kindDecoders      = [numKinds]ValDecoder{}
	keyDecoders       = [numKinds]ValDecoder{}
)

func createGlobalValDecoder(ptr interface{}, dec ValDecoder) {
	ef := (*eface)(unsafe.Pointer(&ptr))
	globalValDecoders[ef.rtype] = dec
}

func mapKind(ptr interface{}, dec ValDecoder) {
	ef := (*eface)(unsafe.Pointer(&ptr))
	kind := reflect.TypeOf(ptr).Elem().Kind()
	kindMap[kind] = ef.rtype
	kindDecoders[kind] = dec
}

func mapKeyDecoder(ptr interface{}, dec ValDecoder) {
	ptrType := reflect.TypeOf(ptr)
	kind := ptrType.Elem().Kind()
	keyDecoders[kind] = dec
}

func init() {
	// standard json library types
	createGlobalValDecoder((*json.Number)(nil), (*jsonNumberDecoder)(nil))
	createGlobalValDecoder((*json.RawMessage)(nil), (*jsonRawMessageDecoder)(nil))

	// kind mapping
	mapKind((*bool)(nil), (*boolDecoder)(nil))
	mapKind((*string)(nil), (*stringDecoder)(nil))
	if strconv.IntSize == 32 {
		mapKind((*int)(nil), (*int32Decoder)(nil))
		mapKind((*uint)(nil), (*uint32Decoder)(nil))
	} else {
		mapKind((*int)(nil), (*int64Decoder)(nil))
		mapKind((*uint)(nil), (*uint64Decoder)(nil))
	}
	if unsafe.Sizeof(uintptr(0)) == 4 {
		mapKind((*uintptr)(nil), (*uint32Decoder)(nil))
	} else {
		mapKind((*uintptr)(nil), (*uint64Decoder)(nil))
	}
	mapKind((*int8)(nil), (*int8Decoder)(nil))
	mapKind((*int16)(nil), (*int16Decoder)(nil))
	mapKind((*int32)(nil), (*int32Decoder)(nil))
	mapKind((*int64)(nil), (*int64Decoder)(nil))
	mapKind((*uint8)(nil), (*uint8Decoder)(nil))
	mapKind((*uint16)(nil), (*uint16Decoder)(nil))
	mapKind((*uint32)(nil), (*uint32Decoder)(nil))
	mapKind((*uint64)(nil), (*uint64Decoder)(nil))
	mapKind((*float32)(nil), (*float32Decoder)(nil))
	mapKind((*float64)(nil), (*float64Decoder)(nil))

	// object key decoders
	mapKeyDecoder((*string)(nil), (*stringKeyDecoder)(nil))
	if strconv.IntSize == 32 {
		mapKeyDecoder((*int)(nil), (*int32KeyDecoder)(nil))
	} else {
		mapKeyDecoder((*int)(nil), (*int64KeyDecoder)(nil))
	}
	mapKeyDecoder((*int8)(nil), (*int8KeyDecoder)(nil))
	mapKeyDecoder((*int16)(nil), (*int16KeyDecoder)(nil))
	mapKeyDecoder((*int32)(nil), (*int32KeyDecoder)(nil))
	mapKeyDecoder((*int64)(nil), (*int64KeyDecoder)(nil))
	mapKeyDecoder((*uint8)(nil), (*uint8KeyDecoder)(nil))
	mapKeyDecoder((*uint16)(nil), (*uint16KeyDecoder)(nil))
	mapKeyDecoder((*uint32)(nil), (*uint32KeyDecoder)(nil))
	mapKeyDecoder((*uint64)(nil), (*uint64KeyDecoder)(nil))
	if unsafe.Sizeof(uintptr(0)) == 4 {
		mapKeyDecoder((*uintptr)(nil), (*uint32KeyDecoder)(nil))
	} else {
		mapKeyDecoder((*uintptr)(nil), (*uint64KeyDecoder)(nil))
	}
}

type ValDecoder interface {
	Decode(ptr unsafe.Pointer, it *Iterator) error
}

type notSupportedDecoder string

func (dec notSupportedDecoder) Decode(ptr unsafe.Pointer, it *Iterator) error {
	return TypeNotSupportedError(dec)
}
