package jzon

import (
	"reflect"
	"sync"
	"sync/atomic"
)

var (
	DefaultEncoder = NewEncoder(nil)
)

type EncoderOption struct {
	ValEncoders map[reflect.Type]ValEncoder

	EscapeHTML bool
	Tag        string
}

type encoderCache = map[rtype]ValEncoder

type Encoder struct {
	cacheMu      sync.Mutex
	encoderCache atomic.Value
	escapeHtml   bool
	tag          string
}

func NewEncoder(opt *EncoderOption) *Encoder {
	enc := Encoder{
		tag:        "json",
		escapeHtml: true,
	}
	cache := encoderCache{}
	if opt != nil {
		// TODO: encoder cache
		enc.escapeHtml = opt.EscapeHTML
		if opt.Tag != "" {
			enc.tag = opt.Tag
		}
	}
	enc.encoderCache.Store(cache)
	return &enc
}
