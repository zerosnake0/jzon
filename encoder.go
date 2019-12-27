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
	cacheMu sync.Mutex
	// the encoder cache, or root encoder cache
	encoderCache atomic.Value
	// the internal cache
	internalCache encoderCache

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
	internalCache := encoderCache{}
	if opt != nil {
		for typ, valEnc := range opt.ValEncoders {
			rtype := rtypeOfType(typ)
			cache[rtype] = valEnc
			internalCache[rtype] = valEnc
		}
		enc.escapeHtml = opt.EscapeHTML
		if opt.Tag != "" {
			enc.tag = opt.Tag
		}
		enc.onlyTaggedField = opt.OnlyTaggedField
	}
	enc.encoderCache.Store(cache)
	enc.internalCache = internalCache
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
	enc.createEncoderInternal(newCache, enc.internalCache, typ)
	enc.encoderCache.Store(newCache)
	return newCache[rtype]
}

func (enc *Encoder) createEncoderInternal(cache, internalCache encoderCache, typesToCreate ...reflect.Type) {
	rebuildMap := map[rtype]interface{}{}
	idx := len(typesToCreate) - 1
	for idx >= 0 {
		typ := typesToCreate[idx]
		typesToCreate = typesToCreate[:idx]
		idx -= 1

		rType := rtypeOfType(typ)
		if _, ok := internalCache[rType]; ok { // check if visited
			continue
		}

		// check global encoders
		if v, ok := globalValEncoders[rType]; ok {
			internalCache[rType] = v
			cache[rType] = v
			continue
		}

		// check json.Marshaler interface
		if typ.Implements(jsonMarshalerType) {
			v := jsonMarshalerEncoder(rType)
			internalCache[rType] = v
			cache[rType] = v
			continue
		}
		// TODO: ptr to json.Marshaler

		// check text.Marshaler interface
		if typ.Implements(textMarshalerType) {
			v := textMarshalerEncoder(rType)
			internalCache[rType] = v
			cache[rType] = v
			continue
		}
		// TODO: ptr to text.Marshaler

		kind := typ.Kind()
		if kindRType := encoderKindMap[kind]; kindRType != 0 {
			// TODO: shall we make this an option?
			// TODO: so that only the native type is affected?
			// check if the native type has a custom encoder
			if v, ok := internalCache[kindRType]; ok {
				internalCache[rType] = v
				cache[rType] = v
				continue
			}

			if v := kindEncoders[kind]; v != nil {
				internalCache[rType] = v
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
			internalCache[rType] = w.encoder
			rebuildMap[rType] = w
		case reflect.Array:
			elemType := typ.Elem()
			typesToCreate = append(typesToCreate, elemType)
			idx += 1
			if typ.Len() == 0 {
				v := (*emptyArrayEncoder)(nil)
				internalCache[rType] = v
				cache[rType] = v
			} else {
				w := newArrayEncoder(typ)
				internalCache[rType] = w.encoder
				rebuildMap[rType] = w
			}
		case reflect.Interface:
			var v ValEncoder
			if typ.NumMethod() == 0 {
				v = (*efaceEncoder)(nil)
			} else {
				v = (*ifaceEncoder)(nil)
			}
			internalCache[rType] = v
			cache[rType] = v
		case reflect.Map:
			w := newMapEncoder(typ)
			if w == nil {
				v := notSupportedEncoder(typ.String())
				internalCache[rType] = v
				cache[rType] = v
			} else {
				typesToCreate = append(typesToCreate, typ.Elem())
				idx += 1
				// pointer decoder is a reverse of direct encoder
				internalCache[rType] = &pointerEncoder{w.encoder}
				rebuildMap[rType] = w
			}
		case reflect.Slice:
			w := newSliceEncoder(typ)
			typesToCreate = append(typesToCreate, typ.Elem())
			idx += 1
			internalCache[rType] = w.encoder
			rebuildMap[rType] = w
		case reflect.Struct:
			fallthrough
		default:
			v := notSupportedEncoder(typ.String())
			internalCache[rType] = v
			cache[rType] = v
		}
	}
	// rebuild some encoders
	for rType, builder := range rebuildMap {
		switch x := builder.(type) {
		case *pointerEncoderBuilder:
			v := internalCache[x.elemRType]
			x.encoder.encoder = v
			cache[rType] = v
		case *arrayEncoderBuilder:
			v := internalCache[x.elemRType]
			x.encoder.encoder = v
			if ifaceIndir(rType) {
				cache[rType] = x.encoder
			} else {
				// (see reflect.ArrayOf)
				// when the array is stored in interface directly, it means:
				// 1. the length of array is 1
				// 2. the element of the array is also directly saved
				cache[rType] = &directEncoder{x.encoder}
			}
		case *mapEncoderBuilder:
			// TODO: key/value encoder
			x.encoder.elemEncoder = internalCache[x.elemRType]
			cache[rType] = x.encoder
		case *sliceEncoderBuilder:
			x.encoder.elemEncoder = internalCache[x.elemRType]
			cache[rType] = x.encoder
		}
	}
}
