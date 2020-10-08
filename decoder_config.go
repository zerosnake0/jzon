package jzon

import (
	"io"
	"reflect"
	"sync"
	"sync/atomic"
)

var (
	// Default decoder config is compatible with standard lib
	DefaultDecoderConfig = NewDecoderConfig(nil)
)

type DecoderOption struct {
	// custom value decoders
	ValDecoders map[reflect.Type]ValDecoder

	// if the object key is case sensitive
	// `false` by default
	CaseSensitive bool

	// the tag name for structures
	// `json` by default
	Tag string

	OnlyTaggedField bool

	UseNumber bool

	DisallowUnknownFields bool
}

type decoderCache = map[rtype]ValDecoder

type DecoderConfig struct {
	cacheMu      sync.Mutex
	decoderCache atomic.Value

	// fixed config, cannot override during runtime
	caseSensitive   bool
	tag             string
	onlyTaggedField bool

	// can override during runtime
	useNumber             bool
	disallowUnknownFields bool
}

func NewDecoderConfig(opt *DecoderOption) *DecoderConfig {
	decCfg := DecoderConfig{
		tag: "json",
	}
	// add decoders to cache
	cache := decoderCache{}
	if opt != nil {
		for elemTyp, valDec := range opt.ValDecoders {
			cache[rtypeOfType(reflect.PtrTo(elemTyp))] = valDec
		}
		decCfg.caseSensitive = opt.CaseSensitive
		if opt.Tag != "" {
			decCfg.tag = opt.Tag
		}
		decCfg.onlyTaggedField = opt.OnlyTaggedField
		decCfg.useNumber = opt.UseNumber
		decCfg.disallowUnknownFields = opt.DisallowUnknownFields
	}
	decCfg.decoderCache.Store(cache)
	return &decCfg
}

func (decCfg *DecoderConfig) Unmarshal(data []byte, obj interface{}) error {
	it := decCfg.NewIterator()
	err := it.Unmarshal(data, obj)
	if err != nil {
		err = it.WrapError(err)
	}
	it.Release()
	return err
}

func (decCfg *DecoderConfig) UnmarshalFromString(s string, obj interface{}) error {
	return decCfg.Unmarshal(localStringToBytes(s), obj)
}

func (decCfg *DecoderConfig) UnmarshalFromReader(r io.Reader, obj interface{}) error {
	it := decCfg.NewIterator()
	err := it.UnmarshalFromReader(r, obj)
	if err != nil {
		err = it.WrapError(err)
	}
	it.Release()
	return err
}

func (decCfg *DecoderConfig) getDecoderFromCache(rType rtype) ValDecoder {
	return decCfg.decoderCache.Load().(decoderCache)[rType]
}

func (decCfg *DecoderConfig) NewDecoder(r io.Reader) *Decoder {
	it := decCfg.NewIterator()
	it.Reset(r)
	return &Decoder{
		it: it,
	}
}

// the typ must be a pointer type
func (decCfg *DecoderConfig) createDecoder(rType rtype, ptrType reflect.Type) ValDecoder {
	decCfg.cacheMu.Lock()
	defer decCfg.cacheMu.Unlock()
	cache := decCfg.decoderCache.Load().(decoderCache)
	// double check
	if vd := cache[rType]; vd != nil {
		return vd
	}
	// make copy
	newCache := decoderCache{}
	for k, v := range cache {
		newCache[k] = v
	}
	var q typeQueue
	q.push(ptrType)
	decCfg.createDecoderInternal(newCache, q)
	decCfg.decoderCache.Store(newCache)
	return newCache[rType]
}

type decoderBuilder interface {
	build(cache decoderCache)
}

func (decCfg *DecoderConfig) createDecoderInternal(cache decoderCache, typesToCreate typeQueue) {
	rebuildMap := map[rtype]decoderBuilder{}
	for ptrType := typesToCreate.pop(); ptrType != nil; ptrType = typesToCreate.pop() {
		rType := rtypeOfType(ptrType)
		if _, ok := cache[rType]; ok { // check if visited
			continue
		}
		// check global decoders
		if v, ok := globalValDecoders[rType]; ok {
			cache[rType] = v
			continue
		}
		// check json.Unmarshaler interface
		if ptrType.Implements(jsonUnmarshalerType) {
			cache[rType] = jsonUnmarshalerDecoder(rType)
			continue
		}
		if ptrType.Implements(textUnmarshalerType) {
			cache[rType] = textUnmarshalerDecoder(rType)
			continue
		}
		elem := ptrType.Elem()
		elemKind := elem.Kind()
		if elemNativeRType := decoderKindMap[elemKind]; elemNativeRType != 0 {
			// TODO: shall we make this an option?
			// TODO: so that only the native type is affected?
			// check if the native type has a custom decoder
			if v, ok := cache[elemNativeRType]; ok {
				cache[rType] = v
				continue
			}
			// otherwise check default native type decoder
			if v := kindDecoders[elemKind]; v != nil {
				cache[rType] = v
				continue
			}
		}
		switch elemKind {
		case reflect.Interface:
			if elem.NumMethod() == 0 {
				cache[rType] = (*efaceDecoder)(nil)
			} else {
				cache[rType] = (*ifaceDecoder)(nil)
			}
		case reflect.Struct:
			fields := describeStruct(elem, decCfg.tag, decCfg.onlyTaggedField)
			numFields := len(fields)
			if numFields == 0 {
				cache[rType] = (*emptyObjectDecoder)(nil)
				continue
			}
			for i := range fields {
				fi := &fields[i]
				typesToCreate.push(fi.ptrType)
			}
			if numFields == 1 {
				w := newOneFieldStructDecoder(&fields[0], decCfg.caseSensitive)
				cache[rType] = w.decoder
				rebuildMap[rType] = w
			} else if numFields <= 10 {
				// TODO: determinate the threshold, several factors may be involved:
				// TODO:   1. number of fields
				// TODO:   2. (average) field length
				// TODO:   3. field similarity
				w := newSmallStructDecoder(fields)
				cache[rType] = w.decoder
				rebuildMap[rType] = w
			} else {
				w := newStructDecoder(fields)
				cache[rType] = w.decoder
				rebuildMap[rType] = w
			}
		case reflect.Ptr:
			typesToCreate.push(elem)
			w := newPointerDecoder(elem)
			cache[rType] = w.decoder
			rebuildMap[rType] = w
		case reflect.Array:
			elemPtrType := reflect.PtrTo(elem.Elem())
			typesToCreate.push(elemPtrType)
			w := newArrayDecoder(elem)
			cache[rType] = w.decoder
			rebuildMap[rType] = w
		case reflect.Slice:
			elemPtrType := reflect.PtrTo(elem.Elem())
			typesToCreate.push(elemPtrType)
			w := newSliceDecoder(elem)
			cache[rType] = w.decoder
			rebuildMap[rType] = w
		case reflect.Map:
			w := newMapDecoder(elem)
			if w == nil {
				cache[rType] = notSupportedDecoder(ptrType.String())
			} else {
				valuePtrType := reflect.PtrTo(elem.Elem())
				typesToCreate.push(valuePtrType)
				cache[rType] = w.decoder
				rebuildMap[rType] = w
			}
		default:
			cache[rType] = notSupportedDecoder(ptrType.String())
		}
	}
	// rebuild some decoders
	for _, builder := range rebuildMap {
		builder.build(cache)
	}
}
