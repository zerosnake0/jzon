package jzon

import (
	"log"
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

// check if the type has a special addr encoder
func (cache encoderCache2) addrEncoder(typ reflect.Type) ValEncoder2 {
	ptrType := reflect.PtrTo(typ)
	ptrEnc, ok := cache[ptrType]
	if !ok {
		return nil
	}
	if _, ok := ptrEnc.(*pointerEncoder2); ok {
		return nil
	}
	return ptrEnc
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

type typQueue2 []reflect.Type

func (tq *typQueue2) push(t reflect.Type) {
	*tq = append(*tq, t)
}

func (tq *typQueue2) pushAlsoPtr(t reflect.Type) {
	if t.Kind() == reflect.Ptr {
		*tq = append(*tq, t)
	} else {
		*tq = append(*tq, reflect.PtrTo(t), t)
	}
}

func (tq *typQueue2) pop() (t reflect.Type) {
	q := *tq
	l := len(q)
	if l == 0 {
		return nil
	}
	t = q[l-1]
	*tq = q[:l-1]
	return
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
	var tq typQueue2
	tq.push(typ)
	enc.createEncoderInternal2(newCache, tq)
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

func (enc *Encoder) createEncoderInternal2(cache encoderCache2, typesToCreate typQueue2) {
	rebuildMap := map[reflect.Type]interface{}{}
	for typ := typesToCreate.pop(); typ != nil; typ = typesToCreate.pop() {
		log.Println("createEncoderInternal2", typ)
		if cache.has(typ) {
			continue
		}
		if _, ok := rebuildMap[typ]; ok { // for some special types
			continue
		}

		if v, ok := globalValEncoders2[typ]; ok {
			cache[typ] = v
			continue
		}

		if typ.Implements(jsonMarshalerType) {
			if typ.Kind() == reflect.Ptr {
				cache[typ] = (*jsonPointerMarshalerEncoder2)(nil)
			} else {
				cache[typ] = (*jsonMarshalerEncoder2)(nil)
			}
			continue
		}

		if typ.Implements(textMarshalerType) {
			if typ.Kind() == reflect.Ptr {
				cache[typ] = (*textPointerMarshalerEncoder2)(nil)
			} else {
				cache[typ] = (*textMarshalerEncoder2)(nil)
			}
			continue
		}

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
			typesToCreate.push(elemType)
			w := newPointerEncoder2(elemType)
			rebuildMap[typ] = w
			cache[typ] = w.encoder
		case reflect.Array:
			if typ.Len() == 0 {
				cache[typ] = (*emptyArrayEncoder)(nil)
			} else {
				typesToCreate.pushAlsoPtr(typ.Elem())
				w := newArrayEncoder2(typ)
				rebuildMap[typ] = w
				cache[typ] = w.encoder
			}
		case reflect.Slice:
			// the elements of a slice must be addressable
			typesToCreate.pushAlsoPtr(typ.Elem())
			w := newSliceEncoder2(typ)
			rebuildMap[typ] = w
			cache[typ] = w.encoder
		case reflect.Interface:
			cache[typ] = (*efaceEncoder)(nil)
		case reflect.Map:
			w := newMapEncoder2(typ)
			if w == nil {
				cache[typ] = notSupportedEncoder(typ.String())
			} else {
				// the map element is never addressable
				typesToCreate.push(typ.Elem())
				rebuildMap[typ] = w
				cache[typ] = w.encoder
			}
		}
	}
	// rebuild base64 encoders
	for typ, builder := range rebuildMap {
		switch x := builder.(type) {
		case *sliceEncoderBuilder2:
			if x.elemType.Kind() != reflect.Uint8 {
				continue
			}
			addrEnc := cache.addrEncoder(x.elemType)
			if addrEnc != nil {
				// the element has a special element addr encoder
				continue
			}
			enc := cache[x.elemType]
			if enc != (*uint8Encoder)(nil) {
				continue
			}
			cache[typ] = (*base64Encoder)(nil)
			delete(rebuildMap, typ)
		}
	}
	// rebuild conditional encoder
	for typ, typEnc := range cache {
		addrEnc := cache.addrEncoder(typ)
		if addrEnc == nil {
			continue
		}
		cache[typ] = &conditionalEncoder{
			addrEnc:  addrEnc,
			valueEnc: typEnc,
		}
	}
	// other encoders
	for _, builder := range rebuildMap {
		switch x := builder.(type) {
		case *pointerEncoderBuilder2:
			v := cache[x.elemType]
			if _, ok := v.(*conditionalEncoder); ok {
				panic("should not reach here")
			}
			x.encoder.encoder = v
		case *arrayEncoderBuilder2:
			x.encoder.encoder = cache[x.elemType]
		case *sliceEncoderBuilder2:
			v := cache[x.elemType]
			if ce, ok := v.(*conditionalEncoder); ok {
				v = &addrEncoder{ce.addrEnc}
			}
			x.encoder.elemEncoder = v
		case *mapEncoderBuilder2:
			v := cache[x.mapType.Elem()]
			if ce, ok := v.(*conditionalEncoder); ok {
				// the map element is never addressable
				v = ce.valueEnc
			}
			x.encoder.elemEncoder = v
		}
	}
}

type conditionalEncoder struct {
	valueEnc ValEncoder2
	addrEnc  ValEncoder2
}

func (ce *conditionalEncoder) Encode2(v reflect.Value, s *Streamer, opts *EncOpts) {
	if v.CanAddr() {
		ce.addrEnc.Encode2(v.Addr(), s, opts)
	} else {
		ce.valueEnc.Encode2(v, s, opts)
	}
}

type addrEncoder struct {
	addrEnc ValEncoder2
}

func (ae *addrEncoder) Encode2(v reflect.Value, s *Streamer, opts *EncOpts) {
	ae.addrEnc.Encode2(v.Addr(), s, opts)
}
