package jzon

import (
	"encoding"
	"encoding/json"
	"reflect"
	"strconv"
	"unsafe"
)

var (
	globalValEncoders = map[rtype]ValEncoder{}
	encoderKindMap    = [numKinds]rtype{}
	kindEncoders      = [numKinds]ValEncoder{}
	keyEncoders       = [numKinds]keyEncoder{}
)

func createGlobalValEncoder(ptr interface{}, enc ValEncoder) {
	typ := reflect.TypeOf(ptr).Elem()
	rType := rtypeOfType(typ)
	if !ifaceIndir(rType) {
		panic("not supported")
	}
	globalValEncoders[rType] = enc
}

func mapEncoderKind(ptr interface{}, enc ValEncoder) {
	elem := reflect.TypeOf(ptr).Elem()
	kind := elem.Kind()
	elemRType := rtypeOfType(elem)
	if !ifaceIndir(elemRType) {
		panic("not supported")
	}
	encoderKindMap[kind] = elemRType
	kindEncoders[kind] = enc
}

func mapKeyEncoder(ptr interface{}, enc keyEncoder) {
	ptrType := reflect.TypeOf(ptr)
	kind := ptrType.Elem().Kind()
	keyEncoders[kind] = enc
}

func init() {
	// standard json library types
	createGlobalValEncoder((*json.Number)(nil), (*jsonNumberEncoder)(nil))
	createGlobalValEncoder((*json.RawMessage)(nil), (*jsonRawMessageEncoder)(nil))
	createGlobalValEncoder((*json.Marshaler)(nil), (*dynamicJsonMarshalerEncoder)(nil))
	createGlobalValEncoder((*encoding.TextMarshaler)(nil), (*dynamicTextMarshalerEncoder)(nil))

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

	// object key encoders
	mapKeyEncoder((*string)(nil), (*stringKeyEncoder)(nil))
	if strconv.IntSize == 32 {
		mapKeyEncoder((*int)(nil), (*int32KeyEncoder)(nil))
		mapKeyEncoder((*uint)(nil), (*uint32KeyEncoder)(nil))
	} else {
		mapKeyEncoder((*int)(nil), (*int64KeyEncoder)(nil))
		mapKeyEncoder((*uint)(nil), (*uint64KeyEncoder)(nil))
	}
	mapKeyEncoder((*int8)(nil), (*int8KeyEncoder)(nil))
	mapKeyEncoder((*int16)(nil), (*int16KeyEncoder)(nil))
	mapKeyEncoder((*int32)(nil), (*int32KeyEncoder)(nil))
	mapKeyEncoder((*int64)(nil), (*int64KeyEncoder)(nil))
	mapKeyEncoder((*uint8)(nil), (*uint8KeyEncoder)(nil))
	mapKeyEncoder((*uint16)(nil), (*uint16KeyEncoder)(nil))
	mapKeyEncoder((*uint32)(nil), (*uint32KeyEncoder)(nil))
	mapKeyEncoder((*uint64)(nil), (*uint64KeyEncoder)(nil))
	if unsafe.Sizeof(uintptr(0)) == 4 {
		mapKeyEncoder((*uintptr)(nil), (*uint32KeyEncoder)(nil))
	} else {
		mapKeyEncoder((*uintptr)(nil), (*uint64KeyEncoder)(nil))
	}
}

type EncOpts struct {
	Quoted bool
}

func (opts *EncOpts) noescape() *EncOpts {
	return (*EncOpts)(noescape(unsafe.Pointer(opts)))
}

type keyEncoder interface {
	Encode(ptr unsafe.Pointer, s *Streamer)
}

type ValEncoder interface {
	IsEmpty(ptr unsafe.Pointer) bool

	Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts)
}

type notSupportedEncoder string

func (enc notSupportedEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func (enc notSupportedEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
	s.Error = TypeNotSupportedError(enc)
}
