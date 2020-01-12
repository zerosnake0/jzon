package jzon

import (
	"encoding/json"
	"reflect"
	"strconv"
	"unsafe"
)

var (
	globalValDecoders = map[rtype]ValDecoder{}
	decoderKindMap    = [numKinds]rtype{}
	kindDecoders      = [numKinds]ValDecoder{}
	keyDecoders       = [numKinds]ValDecoder{}
)

func createGlobalValDecoder(ptr interface{}, dec ValDecoder) {
	ef := (*eface)(unsafe.Pointer(&ptr))
	globalValDecoders[ef.rtype] = dec
}

func mapDecoderKind(ptr interface{}, dec ValDecoder) {
	ef := (*eface)(unsafe.Pointer(&ptr))
	kind := reflect.TypeOf(ptr).Elem().Kind()
	decoderKindMap[kind] = ef.rtype
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
	mapDecoderKind((*bool)(nil), (*boolDecoder)(nil))
	mapDecoderKind((*string)(nil), (*stringDecoder)(nil))
	if strconv.IntSize == 32 {
		mapDecoderKind((*int)(nil), (*int32Decoder)(nil))
		mapDecoderKind((*uint)(nil), (*uint32Decoder)(nil))
	} else {
		mapDecoderKind((*int)(nil), (*int64Decoder)(nil))
		mapDecoderKind((*uint)(nil), (*uint64Decoder)(nil))
	}
	if unsafe.Sizeof(uintptr(0)) == 4 {
		mapDecoderKind((*uintptr)(nil), (*uint32Decoder)(nil))
	} else {
		mapDecoderKind((*uintptr)(nil), (*uint64Decoder)(nil))
	}
	mapDecoderKind((*int8)(nil), (*int8Decoder)(nil))
	mapDecoderKind((*int16)(nil), (*int16Decoder)(nil))
	mapDecoderKind((*int32)(nil), (*int32Decoder)(nil))
	mapDecoderKind((*int64)(nil), (*int64Decoder)(nil))
	mapDecoderKind((*uint8)(nil), (*uint8Decoder)(nil))
	mapDecoderKind((*uint16)(nil), (*uint16Decoder)(nil))
	mapDecoderKind((*uint32)(nil), (*uint32Decoder)(nil))
	mapDecoderKind((*uint64)(nil), (*uint64Decoder)(nil))
	mapDecoderKind((*float32)(nil), (*float32Decoder)(nil))
	mapDecoderKind((*float64)(nil), (*float64Decoder)(nil))

	// object key decoders
	mapKeyDecoder((*string)(nil), (*stringDecoder)(nil))
	if strconv.IntSize == 32 {
		mapKeyDecoder((*int)(nil), (*int32Decoder)(nil))
		mapKeyDecoder((*uint)(nil), (*uint32Decoder)(nil))
	} else {
		mapKeyDecoder((*int)(nil), (*int64Decoder)(nil))
		mapKeyDecoder((*uint)(nil), (*uint64Decoder)(nil))
	}
	mapKeyDecoder((*int8)(nil), (*int8Decoder)(nil))
	mapKeyDecoder((*int16)(nil), (*int16Decoder)(nil))
	mapKeyDecoder((*int32)(nil), (*int32Decoder)(nil))
	mapKeyDecoder((*int64)(nil), (*int64Decoder)(nil))
	mapKeyDecoder((*uint8)(nil), (*uint8Decoder)(nil))
	mapKeyDecoder((*uint16)(nil), (*uint16Decoder)(nil))
	mapKeyDecoder((*uint32)(nil), (*uint32Decoder)(nil))
	mapKeyDecoder((*uint64)(nil), (*uint64Decoder)(nil))
	if unsafe.Sizeof(uintptr(0)) == 4 {
		mapKeyDecoder((*uintptr)(nil), (*uint32Decoder)(nil))
	} else {
		mapKeyDecoder((*uintptr)(nil), (*uint64Decoder)(nil))
	}
}

type DecOpts struct {
	MapKey bool
	Quoted bool
}

func (opts *DecOpts) noescape() *DecOpts {
	return (*DecOpts)(noescape(unsafe.Pointer(opts)))
}

type ValDecoder interface {
	Decode(ptr unsafe.Pointer, it *Iterator, opts *DecOpts) error
}

type notSupportedDecoder string

func (dec notSupportedDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	return TypeNotSupportedError(dec)
}
