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

type encoderCache map[rtype]ValEncoder

func (cache encoderCache) has(rtype rtype) bool {
	_, ok := cache[rtype]
	return ok
}

// make sure that the pointer encoders has already been rebuilt
// before calling, so it's safe to use it's internal encoder
func (cache encoderCache) preferPtrEncoder(typ reflect.Type) ValEncoder {
	ptrType := reflect.PtrTo(typ)
	ptrEncoder := cache[rtypeOfType(ptrType)]
	if pe, ok := ptrEncoder.(*pointerEncoder); ok {
		return pe.encoder
	} else {
		// the element has a special pointer encoder
		return &directEncoder{ptrEncoder}
	}
}

type encoderCache2 map[reflect.Type]ValEncoder2

func (cache encoderCache2) has(typ reflect.Type) bool {
	_, ok := cache[typ]
	return ok
}

type Encoder struct {
	cacheMu sync.Mutex
	// the encoder cache, or root encoder cache
	encoderCache atomic.Value
	// the internal cache
	internalCache encoderCache

	encoderCache2 atomic.Value

	escapeHtml      bool
	safeSet         []string
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
	cache2 := encoderCache2{}
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
	enc.encoderCache2.Store(cache2)
	if enc.escapeHtml {
		enc.safeSet = htmlSafeSet[:]
	} else {
		enc.safeSet = safeSet[:]
	}
	return &enc
}

func (enc *Encoder) Marshal(obj interface{}) ([]byte, error) {
	s := enc.NewStreamer()
	defer enc.ReturnStreamer(s)
	s.Value(obj)
	if s.Error != nil {
		return nil, s.Error
	}
	// we make a new slice with explicit size,
	//   1. the internal buffer may be much longer than the output one,
	//      it can be used for longer output
	//   2. avoid calling bytes buffer pool (sync.Pool)
	b := make([]byte, len(s.buffer))
	copy(b, s.buffer)
	return b, nil
}

func (enc *Encoder) Marshal2(obj interface{}) ([]byte, error) {
	s := enc.NewStreamer()
	defer enc.ReturnStreamer(s)
	s.Value2(obj)
	if s.Error != nil {
		return nil, s.Error
	}
	// we make a new slice with explicit size,
	//   1. the internal buffer may be much longer than the output one,
	//      it can be used for longer output
	//   2. avoid calling bytes buffer pool (sync.Pool)
	b := make([]byte, len(s.buffer))
	copy(b, s.buffer)
	return b, nil
}

func (enc *Encoder) getEncoderFromCache(rtype rtype) ValEncoder {
	return enc.encoderCache.Load().(encoderCache)[rtype]
}

func (enc *Encoder) getEncoderFromCache2(typ reflect.Type) ValEncoder2 {
	return enc.encoderCache2.Load().(encoderCache2)[typ]
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

func (enc *Encoder) createEncoder2(typ reflect.Type) ValEncoder2 {
	enc.cacheMu.Lock()
	defer enc.cacheMu.Unlock()
	cache := enc.encoderCache2.Load().(encoderCache2)
	if ve := cache[typ]; ve != nil {
		return ve
	}
	newCache := encoderCache2{}
	for k, v := range cache {
		newCache[k] = v
	}
	enc.createEncoderInternal2(newCache, typ)
	enc.encoderCache2.Store(newCache)
	return newCache[typ]
}

func (enc *Encoder) createEncoderInternal(cache, internalCache encoderCache, typesToCreate ...reflect.Type) {
	rebuildMap := map[rtype]interface{}{}
	idx := len(typesToCreate) - 1
	for idx >= 0 {
		typ := typesToCreate[idx]

		typesToCreate = typesToCreate[:idx]
		idx -= 1

		rType := rtypeOfType(typ)
		if internalCache.has(rType) { // check if visited
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
			if ifaceIndir(rType) {
				v := jsonMarshalerEncoder(rType)
				internalCache[rType] = v
				cache[rType] = v
				continue
			}
			if typ.Kind() == reflect.Ptr {
				elemType := typ.Elem()
				if elemType.Implements(jsonMarshalerType) {
					// treat as a pointer encoder
					typesToCreate = append(typesToCreate, elemType)
					idx += 1
					w := newPointerEncoder(elemType)
					internalCache[rType] = w.encoder
					rebuildMap[rType] = w
					continue
				}
			}
			v := directJsonMarshalerEncoder(rType)
			internalCache[rType] = v
			cache[rType] = &directEncoder{v}
			continue
		}

		// check encoding.TextMarshaler interface
		if typ.Implements(textMarshalerType) {
			if ifaceIndir(rType) {
				v := textMarshalerEncoder(rType)
				internalCache[rType] = v
				cache[rType] = v
				continue
			}
			if typ.Kind() == reflect.Ptr {
				elemType := typ.Elem()
				if elemType.Implements(textMarshalerType) {
					// treat as a pointer encoder
					typesToCreate = append(typesToCreate, elemType)
					idx += 1
					w := newPointerEncoder(elemType)
					internalCache[rType] = w.encoder
					rebuildMap[rType] = w
					continue
				}
			}
			v := directTextMarshalerEncoder(rType)
			internalCache[rType] = v
			cache[rType] = &directEncoder{v}
			continue
		}

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
			typesToCreate = append(typesToCreate, reflect.PtrTo(typ.Elem()))
			idx += 1
			internalCache[rType] = w.encoder
			rebuildMap[rType] = w
		case reflect.Struct:
			w := enc.newStructEncoder(typ)
			if w == nil {
				// no fields to marshal
				v := (*emptyStructEncoder)(nil)
				internalCache[rType] = v
				cache[rType] = v
			} else {
				for i := range w.fields.list {
					fi := &w.fields.list[i]
					typesToCreate = append(typesToCreate, fi.ptrType)
					idx += 1
				}
				internalCache[rType] = w.encoder
				rebuildMap[rType] = w
			}
		default:
			v := notSupportedEncoder(typ.String())
			internalCache[rType] = v
			cache[rType] = v
		}
	}
	// rebuild base64 encoders
	for rType, builder := range rebuildMap {
		switch x := builder.(type) {
		case *sliceEncoderBuilder:
			if x.elemType.Kind() != reflect.Uint8 {
				continue
			}
			elemPtrType := reflect.PtrTo(x.elemType)
			elemPtrEncoder := internalCache[rtypeOfType(elemPtrType)]
			if _, ok := elemPtrEncoder.(*pointerEncoder); !ok {
				// the element has a special pointer encoder
				continue
			}
			// the pointer decoder has not been rebuilt yet
			// we need to use the explicit element rtype
			elemEncoder := internalCache[rtypeOfType(x.elemType)]
			if elemEncoder != (*uint8Encoder)(nil) {
				// the element has a special value encoder
				continue
			}
			v := (*base64Encoder)(nil)
			internalCache[rType] = v
			cache[rType] = v
			delete(rebuildMap, rType)
		}
	}
	// rebuild ptr encoders
	for rType, builder := range rebuildMap {
		switch x := builder.(type) {
		case *pointerEncoderBuilder:
			v := internalCache[x.elemRType]
			x.encoder.encoder = v
			cache[rType] = v
			delete(rebuildMap, rType)
		}
	}
	// rebuild other encoders
	for rType, builder := range rebuildMap {
		switch x := builder.(type) {
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
			x.encoder.elemEncoder = internalCache.preferPtrEncoder(x.elemType)
			cache[rType] = x.encoder
		case *structEncoderBuilder:
			x.encoder.fields.init(len(x.fields.list))
			for i := range x.fields.list {
				fi := &x.fields.list[i]
				v := internalCache.preferPtrEncoder(fi.ptrType.Elem())
				x.encoder.fields.add(fi, enc.escapeHtml, v)
			}
			if ifaceIndir(rType) {
				cache[rType] = x.encoder
			} else {
				cache[rType] = &directEncoder{x.encoder}
			}
		}
	}
}

func (enc *Encoder) createEncoderInternal2(cache encoderCache2, typesToCreate ...reflect.Type) {
	rebuildMap := map[reflect.Type]interface{}{}
	idx := len(typesToCreate) - 1
	for idx >= 0 {
		typ := typesToCreate[idx]

		typesToCreate = typesToCreate[:idx]
		idx -= 1

		if cache.has(typ) {
			continue
		}

		if v, ok := globalValEncoders2[typ]; ok {
			cache[typ] = v
			continue
		}

		if typ.Implements(jsonMarshalerType) {
			if typ.Kind() == reflect.Ptr {
				cache[typ] = &jsonMarshalerPointerEncoder2{typ}
			} else {
				cache[typ] = &jsonMarshalerEncoder2{typ}
			}
			continue
		}

		// TODO: text marshaler

		kind := typ.Kind()
		if kindType := encoderKindMap2[kind]; kindType != nil {
			if v, ok := cache[kindType]; ok {
				cache[typ] = v
				continue
			}

			if v := kindEncoders2[kind]; v != nil {
				cache[typ] = v
				continue
			}
		}

		switch kind {
		case reflect.Ptr:
			elemType := typ.Elem()
			typesToCreate = append(typesToCreate, elemType)
			idx += 1
			w := newPointerEncoder2(elemType)
			rebuildMap[typ] = w
			cache[typ] = w.encoder
		case reflect.Array:
			if typ.Len() == 0 {
				cache[typ] = (*emptyArrayEncoder)(nil)
			} else {
				elemType := typ.Elem()
				typesToCreate = append(typesToCreate, elemType)
				idx += 1
				w := newArrayEncoder2(typ)
				rebuildMap[typ] = w
				cache[typ] = w.encoder
			}
		case reflect.Slice:
			elemType := typ.Elem()
			if elemType.Kind() == reflect.Uint8 {
				v, ok := cache[elemType]
				if !ok || v == (*uint8Encoder)(nil) {
					cache[typ] = (*base64Encoder)(nil)
					continue
				}
			}
			typesToCreate = append(typesToCreate, elemType)
			idx += 1
			w := newSliceEncoder2(typ)
			rebuildMap[typ] = w
			cache[typ] = w.encoder
		}
	}
	for _, builder := range rebuildMap {
		switch x := builder.(type) {
		case *pointerEncoderBuilder2:
			x.encoder.encoder = cache[x.elemType]
		case *arrayEncoderBuilder2:
			x.encoder.encoder = cache[x.elemType]
		case *sliceEncoderBuilder2:
			x.encoder.elemEncoder = cache[x.elemType]
		}
	}
}
