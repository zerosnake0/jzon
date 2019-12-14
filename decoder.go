package jzon

import (
	"reflect"
	"sync"
	"sync/atomic"
)

var (
	DefaultDecoder = NewDecoder(nil)
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
}

type decoderCache = map[rtype]ValDecoder

type Decoder struct {
	cacheMu      sync.Mutex
	decoderCache atomic.Value

	caseSensitive bool
	tag           string
}

func NewDecoder(opt *DecoderOption) *Decoder {
	dec := Decoder{
		tag: "json",
	}
	// add decoders to cache
	cache := decoderCache{}
	if opt != nil {
		for elemTyp, valDec := range opt.ValDecoders {
			cache[rtypeOfType(reflect.PtrTo(elemTyp))] = valDec
		}
		dec.caseSensitive = opt.CaseSensitive
		if opt.Tag != "" {
			dec.tag = opt.Tag
		}
	}
	dec.decoderCache.Store(cache)
	return &dec
}

func (dec *Decoder) Unmarshal(data []byte, obj interface{}) error {
	it := dec.NewIterator()
	err := it.Unmarshal(data, obj)
	dec.ReturnIterator(it)
	return err
}

func (dec *Decoder) getDecoderFromCache(rType rtype) ValDecoder {
	return dec.decoderCache.Load().(decoderCache)[rType]
}

// the typ must be a pointer type
func (dec *Decoder) createDecoder(rType rtype, ptrType reflect.Type) ValDecoder {
	dec.cacheMu.Lock()
	defer dec.cacheMu.Unlock()
	cache := dec.decoderCache.Load().(decoderCache)
	// double check
	if vd := cache[rType]; vd != nil {
		return vd
	}
	// make copy
	newCache := decoderCache{}
	for k, v := range cache {
		newCache[k] = v
	}
	typesToCreate := []reflect.Type{ptrType}
	dec.createDecoderInternal(newCache, typesToCreate)
	dec.decoderCache.Store(newCache)
	return newCache[rType]
}

func (dec *Decoder) createDecoderInternal(cache decoderCache, typesToCreate []reflect.Type) {
	rebuildMap := decoderCache{}
	idx := len(typesToCreate) - 1
	for idx >= 0 {
		// pop one
		ptrType := typesToCreate[idx]
		typesToCreate = typesToCreate[:idx]
		idx -= 1

		rType := rtypeOfType(ptrType)
		if _, ok := cache[rType]; ok { // double check
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
		if elemNativeRType := kindMap[elemKind]; elemNativeRType != 0 {
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
			vd := dec.newStructDecoder(elem)
			if vd == nil {
				// no field to unmarshal
				cache[rType] = (*skipDecoder)(nil)
			} else {
				for _, fi := range vd.fields {
					typesToCreate = append(typesToCreate, fi.ptrType)
					idx += 1
				}
				cache[rType] = vd
				rebuildMap[rType] = vd
			}
		case reflect.Ptr:
			typesToCreate = append(typesToCreate, elem)
			idx += 1
			vd := newPointerDecoder(elem)
			cache[rType] = vd
			rebuildMap[rType] = vd
		case reflect.Array:
			elemPtrType := reflect.PtrTo(elem.Elem())
			typesToCreate = append(typesToCreate, elemPtrType)
			idx += 1
			vd := newArrayDecoder(elem)
			cache[rType] = vd
			rebuildMap[rType] = vd
		case reflect.Slice:
			elemPtrType := reflect.PtrTo(elem.Elem())
			typesToCreate = append(typesToCreate, elemPtrType)
			idx += 1
			vd := newSliceDecoder(elem)
			cache[rType] = vd
			rebuildMap[rType] = vd
		case reflect.Map:
			vd := newMapDecoder(elem)
			if vd == nil {
				cache[rType] = notSupportedDecoder(ptrType.String())
			} else {
				valuePtrType := reflect.PtrTo(elem.Elem())
				typesToCreate = append(typesToCreate, valuePtrType)
				idx += 1
				cache[rType] = vd
				rebuildMap[rType] = vd
			}
		default:
			cache[rType] = notSupportedDecoder(ptrType.String())
		}
	}
	// rebuild some decoders
	for _, vd := range rebuildMap {
		switch x := vd.(type) {
		case *pointerDecoder:
			x.elemDec = cache[x.ptrRType]
		case *structDecoder:
			for _, field := range x.fields {
				field.decoder = cache[field.rtype]
			}
		case *arrayDecoder:
			x.elemDec = cache[x.elemPtrRType]
		case *sliceDecoder:
			x.elemDec = cache[x.elemPtrRType]
		case *mapDecoder:
			x.valDec = cache[x.valPtrRType]
		}
	}
}
