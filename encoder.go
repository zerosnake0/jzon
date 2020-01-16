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

type Encoder struct {
	cacheMu sync.Mutex
	// the encoder cache, or root encoder cache
	encoderCache atomic.Value
	// the internal cache
	internalCache encoderCache

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
	var q typeQueue
	q.push(typ)
	enc.createEncoderInternal(newCache, enc.internalCache, q)
	enc.encoderCache.Store(newCache)
	return newCache[rtype]
}

func (enc *Encoder) createEncoderInternal(cache, internalCache encoderCache, typesToCreate typeQueue) {
	rebuildMap := map[rtype]interface{}{}
	for typ := typesToCreate.pop(); typ != nil; typ = typesToCreate.pop() {
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

		kind := typ.Kind()

		// check json.Marshaler interface
		if typ.Implements(jsonMarshalerType) {
			if ifaceIndir(rType) {
				v := &jsonMarshalerEncoder{
					isEmpty: isEmptyFunctions[kind],
					rtype:   rType,
				}
				internalCache[rType] = v
				cache[rType] = v
				continue
			}
			if typ.Kind() == reflect.Ptr {
				elemType := typ.Elem()
				if elemType.Implements(jsonMarshalerType) {
					// treat as a pointer encoder
					typesToCreate.push(elemType)
					w := newPointerEncoder(elemType)
					internalCache[rType] = w.encoder
					rebuildMap[rType] = w
				} else {
					v := pointerJsonMarshalerEncoder(rType)
					internalCache[rType] = v
					cache[rType] = &directEncoder{v}
				}
				continue
			}
			v := &directJsonMarshalerEncoder{
				isEmpty: isEmptyFunctions[kind],
				rtype:   rType,
			}
			internalCache[rType] = v
			cache[rType] = &directEncoder{v}
			continue
		}

		// check encoding.TextMarshaler interface
		if typ.Implements(textMarshalerType) {
			if ifaceIndir(rType) {
				v := &textMarshalerEncoder{
					isEmpty: isEmptyFunctions[kind],
					rtype:   rType,
				}
				internalCache[rType] = v
				cache[rType] = v
				continue
			}
			if typ.Kind() == reflect.Ptr {
				elemType := typ.Elem()
				if elemType.Implements(textMarshalerType) {
					// treat as a pointer encoder
					typesToCreate.push(elemType)
					w := newPointerEncoder(elemType)
					internalCache[rType] = w.encoder
					rebuildMap[rType] = w
				} else {
					v := pointerTextMarshalerEncoder(rType)
					internalCache[rType] = v
					cache[rType] = &directEncoder{v}
				}
				continue
			}
			v := &directTextMarshalerEncoder{
				isEmpty: isEmptyFunctions[kind],
				rtype:   rType,
			}
			internalCache[rType] = v
			cache[rType] = &directEncoder{v}
			continue
		}

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
			typesToCreate.push(elemType)
			w := newPointerEncoder(elemType)
			internalCache[rType] = w.encoder
			rebuildMap[rType] = w
		case reflect.Array:
			elemType := typ.Elem()
			typesToCreate.push(reflect.PtrTo(elemType))
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
				typesToCreate.push(typ.Elem())
				// pointer decoder is a reverse of direct encoder
				internalCache[rType] = w.encoder
				rebuildMap[rType] = w
			}
		case reflect.Slice:
			w := newSliceEncoder(typ)
			typesToCreate.push(reflect.PtrTo(typ.Elem()))
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
				for i := range w.fields {
					fi := &w.fields[i]
					typesToCreate.push(fi.ptrType)
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
			x.encoder.encoder = internalCache.preferPtrEncoder(x.elemType)
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
			// TODO: key encoder
			x.encoder.elemEncoder = internalCache[x.elemRType]
			cache[rType] = &directEncoder{x.encoder}
		case *sliceEncoderBuilder:
			x.encoder.elemEncoder = internalCache.preferPtrEncoder(x.elemType)
			cache[rType] = x.encoder
		case *structEncoderBuilder:
			x.encoder.fields.init(len(x.fields))
			for i := range x.fields {
				fi := &x.fields[i]
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
