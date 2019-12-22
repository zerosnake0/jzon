package jzon

import (
	"encoding/json"
	"reflect"
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
}

type ValEncoder interface {
	Encode(ptr unsafe.Pointer, s *Streamer)
}
