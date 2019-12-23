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

	EscapeHTML      bool
	Tag             string
	OnlyTaggedField bool
}

type encoderCache = map[rtype]ValEncoder

type Encoder struct {
	cacheMu      sync.Mutex
	encoderCache atomic.Value

	escapeHtml      bool
	tag             string
	onlyTaggedField bool
}

func NewEncoder(opt *EncoderOption) *Encoder {
	enc := Encoder{
		tag:        "json",
		escapeHtml: true,
	}
	cache := encoderCache{}
	if opt != nil {
		for typ, valEnc := range opt.ValEncoders {
			cache[rtypeOfType(typ)] = valEnc
		}
		enc.escapeHtml = opt.EscapeHTML
		if opt.Tag != "" {
			enc.tag = opt.Tag
		}
		enc.onlyTaggedField = opt.OnlyTaggedField
	}
	enc.encoderCache.Store(cache)
	return &enc
}

func (enc *Encoder) getEncoderFromCache(rtype rtype) ValEncoder {
	return enc.encoderCache.Load().(encoderCache)[rtype]
}

func (enc *Encoder) createEncoder(rtype rtype, typ reflect.Type) ValEncoder {
	enc.cacheMu.Lock()
	defer enc.cacheMu.Unlock()
	cache := enc.encoderCache.Load().(encoderCache)
	// double check
	if ve := cache[rtype]; ve != nil {
		return ve
	}
	newCache := encoderCache{}
	for k, v := range cache {
		newCache[k] = v
	}
	typesToCreate := []reflect.Type{typ}
	enc.createEncoderInternal(newCache, typesToCreate)
	enc.encoderCache.Store(newCache)
	return newCache[rtype]
}

func (enc *Encoder) createEncoderInternal(cache encoderCache, typesToCreate []reflect.Type) {
	rebuildMap := map[rtype]interface{}{}
	idx := len(typesToCreate) - 1
	for idx >= 0 {
		typ := typesToCreate[idx]
		typesToCreate = typesToCreate[:idx]
		idx -= 1

		rType := rtypeOfType(typ)
		if _, ok := cache[rType]; ok { // check if visited
			continue
		}

		// check global encoders
		if v, ok := globalValEncoders[rType]; ok {
			cache[rType] = v
			continue
		}

		// check json.Marshaler interface
		if typ.Implements(jsonMarshalerType) {
			cache[rType] = jsonMarshalerEncoder(rType)
			continue
		}
		// TODO: ptr to json.Marshaler

		// check text.Marshaler interface
		if typ.Implements(textMarshalerType) {
			cache[rType] = textMarshalerEncoder(rType)
			continue
		}
		// TODO: ptr to text.Marshaler

		kind := typ.Kind()
		if kindRType := encoderKindMap[kind]; kindRType != 0 {
			// TODO: shall we make this an option?
			// TODO: so that only the native type is affected?
			// check if the native type has a custom encoder
			if v, ok := cache[kindRType]; ok {
				cache[rType] = v
				continue
			}

			if v := kindEncoders[kind]; v != nil {
				cache[rType] = v
				continue
			}
		}

		switch kind {
		case reflect.Ptr:
			elemType := typ.Elem()
			typesToCreate = append(typesToCreate, elemType)
			idx += 1
			w := newPointerEncoder(elemType)
			cache[rType] = w.encoder
			rebuildMap[rType] = w
		}
	}
	// rebuild some encoders
	for _, builder := range rebuildMap {
		switch x := builder.(type) {
		case *pointerEncoderBuilder:
			x.encoder.encoder = cache[x.rtype]
		}
	}
}
