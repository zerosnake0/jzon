package jzon

import (
	"encoding/json"
	"reflect"
	"unsafe"
)

var (
	globalValEncoders = map[rtype]ValEncoder{}
)

func createGlobalValEncoder(ptr interface{}, enc ValEncoder) {
	typ := reflect.TypeOf(ptr).Elem()
	globalValEncoders[rtypeOfType(typ)] = enc
}

func init() {
	// standard json library types
	createGlobalValEncoder((*json.Number)(nil), (*jsonNumberEncoder)(nil))
}

type ValEncoder interface {
	Encode(ptr unsafe.Pointer, s *Streamer)
}
