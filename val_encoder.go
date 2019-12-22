package jzon

import (
	"encoding/json"
	"reflect"
	"strconv"
	"unsafe"
)

var (
	globalValEncoders = map[rtype]ValEncoder{}
	encoderKindMap    = [numKinds]rtype{}
	kindEncoders      = [numKinds]ValEncoder{}
)

func createGlobalValEncoder(ptr interface{}, enc ValEncoder) {
	typ := reflect.TypeOf(ptr).Elem()
	globalValEncoders[rtypeOfType(typ)] = enc
}

func mapEncoderKind(ptr interface{}, enc ValEncoder) {
	elem := reflect.TypeOf(ptr).Elem()
	kind := elem.Kind()
	encoderKindMap[kind] = rtypeOfType(elem)
	kindEncoders[kind] = enc
}

func init() {
	// standard json library types
	createGlobalValEncoder((*json.Number)(nil), (*jsonNumberEncoder)(nil))
	createGlobalValEncoder((*json.RawMessage)(nil), (*jsonRawMessageEncoder)(nil))

	// kind mapping
	mapEncoderKind((*bool)(nil), (*boolEncoder)(nil))
	mapEncoderKind((*string)(nil), (*stringEncoder)(nil))
	if strconv.IntSize == 32 {
		mapEncoderKind((*int)(nil), (*int32Encoder)(nil))
		mapEncoderKind((*uint)(nil), (*uint32Encoder)(nil))
	} else {
		mapEncoderKind((*int)(nil), (*int64Encoder)(nil))
		mapEncoderKind((*uint)(nil), (*uint64Encoder)(nil))
	}
	if unsafe.Sizeof(uintptr(0)) == 4 {
		mapEncoderKind((*uintptr)(nil), (*uint32Encoder)(nil))
	} else {
		mapEncoderKind((*uintptr)(nil), (*uint64Encoder)(nil))
	}
	mapEncoderKind((*int8)(nil), (*int8Encoder)(nil))
	mapEncoderKind((*int16)(nil), (*int16Encoder)(nil))
	mapEncoderKind((*int32)(nil), (*int32Encoder)(nil))
	mapEncoderKind((*int64)(nil), (*int64Encoder)(nil))
	mapEncoderKind((*uint8)(nil), (*uint8Encoder)(nil))
	mapEncoderKind((*uint16)(nil), (*uint16Encoder)(nil))
	mapEncoderKind((*uint32)(nil), (*uint32Encoder)(nil))
	mapEncoderKind((*uint64)(nil), (*uint64Encoder)(nil))
	mapEncoderKind((*float32)(nil), (*float32Encoder)(nil))
	mapEncoderKind((*float64)(nil), (*float64Encoder)(nil))
}

type ValEncoder interface {
	Encode(ptr unsafe.Pointer, s *Streamer)
}
