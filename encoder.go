package jzon

import (
	"reflect"
	"sync"
	"sync/atomic"
)

var (
	DefaultEncoder = NewEncoder(&EncoderOption{
		//
	})
)

type EncoderOption struct {
	ValEncoders map[reflect.Type]ValEncoder

	Tag string
}

type encoderCache = map[rtype]ValEncoder

type Encoder struct {
	cacheMu      sync.Mutex
	encoderCache atomic.Value
	tag          string
}

func NewEncoder(opt *EncoderOption) *Encoder {
	var enc Encoder
	// TODO:
	if enc.tag == "" {
		enc.tag = "json"
	}
	return &enc
}
