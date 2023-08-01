package jzon

import (
	"io"
	"reflect"
	"sync"
	"sync/atomic"
)

var (
	// DefaultEncoderConfig is compatible with standard lib
	DefaultEncoderConfig = NewEncoderConfig(nil)
)

// EncoderOption can be used to customize the encoder config
type EncoderOption struct {
	ValEncoders   map[reflect.Type]ValEncoder
	IfaceEncoders []IfaceValEncoderConfig

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
	}
	// the element has a special pointer encoder
	return &directEncoder{ptrEncoder}
}

// EncoderConfig is a frozen config for encoding
type EncoderConfig struct {
	cacheMu sync.Mutex
	// the encoder cache, or root encoder cache
	encoderCache atomic.Value
	// the internal cache
	internalCache encoderCache
	// iface encoders
	ifaceEncoderMap map[rtype]ValEncoder
	ifaceEncoder    []IfaceValEncoderConfig

	tag             string
	onlyTaggedField bool

	// can override during runtime
	escapeHTML bool
}

func (encCfg *EncoderConfig) addIfaceEncoder(cfg IfaceValEncoderConfig) {
	rt := rtypeOfType(cfg.Type)
	if encCfg.ifaceEncoderMap[rt] != nil {
		return
	}

	// get the pointer type
	ptrRt := rtypeOfType(reflect.New(cfg.Type).Type())
	encCfg.ifaceEncoderMap[rt] = &dynamicIfaceValEncoder{
		rtype:   ptrRt,
		encoder: cfg.Encoder,
	}
	encCfg.ifaceEncoder = append(encCfg.ifaceEncoder, cfg)
}

// NewEncoderConfig returns a new encoder config
// If the input option is nil, the default option will be applied
func NewEncoderConfig(opt *EncoderOption) *EncoderConfig {
	encCfg := EncoderConfig{
		tag:             "json",
		escapeHTML:      true,
		ifaceEncoderMap: map[rtype]ValEncoder{},
	}
	cache := encoderCache{}
	internalCache := encoderCache{}
	if opt != nil {
		for typ, valEnc := range opt.ValEncoders {
			rtype := rtypeOfType(typ)
			cache[rtype] = valEnc
			internalCache[rtype] = valEnc
		}
		encCfg.escapeHTML = opt.EscapeHTML
		if opt.Tag != "" {
			encCfg.tag = opt.Tag
		}
		encCfg.onlyTaggedField = opt.OnlyTaggedField

		// iface
		encCfg.ifaceEncoder = make([]IfaceValEncoderConfig, 0, len(opt.IfaceEncoders))
		for _, enc := range opt.IfaceEncoders {
			encCfg.addIfaceEncoder(enc)
		}
	}
	encCfg.addIfaceEncoder(IfaceValEncoderConfig{
		Type:    jsonMarshalerType,
		Encoder: jsonMarshalerValEncoder{},
	})
	encCfg.addIfaceEncoder(IfaceValEncoderConfig{
		Type:    textMarshalerType,
		Encoder: textMarshalerValEncoder{},
	})

	encCfg.encoderCache.Store(cache)
	encCfg.internalCache = internalCache
	return &encCfg
}

// Marshal behave like json.Marshal
func (encCfg *EncoderConfig) Marshal(obj interface{}) ([]byte, error) {
	s := encCfg.NewStreamer()
	defer s.Release()
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

// NewEncoder returns a new encoder that writes to w.
func (encCfg *EncoderConfig) NewEncoder(w io.Writer) *Encoder {
	s := encCfg.NewStreamer()
	s.Reset(w)
	return &Encoder{
		s: s,
	}
}

func (encCfg *EncoderConfig) getEncoderFromCache(rtype rtype) ValEncoder {
	return encCfg.encoderCache.Load().(encoderCache)[rtype]
}

func (encCfg *EncoderConfig) createEncoder(rtype rtype, typ reflect.Type) ValEncoder {
	encCfg.cacheMu.Lock()
	defer encCfg.cacheMu.Unlock()
	cache := encCfg.encoderCache.Load().(encoderCache)
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
	encCfg.createEncoderInternal(newCache, encCfg.internalCache, q)
	encCfg.encoderCache.Store(newCache)
	return newCache[rtype]
}

func (encCfg *EncoderConfig) createEncoderInternal(cache, internalCache encoderCache, typesToCreate typeQueue) {
	rebuildMap := map[rtype]interface{}{}
OuterLoop:
	for typ := typesToCreate.pop(); typ != nil; typ = typesToCreate.pop() {
		rType := rtypeOfType(typ)
		if internalCache.has(rType) { // check if visited
			continue
		}

		// check local encoders
		if v, ok := encCfg.ifaceEncoderMap[rType]; ok {
			internalCache[rType] = v
			cache[rType] = v
			continue
		}

		// check global encoders
		if v, ok := globalValEncoders[rType]; ok {
			internalCache[rType] = v
			cache[rType] = v
			continue
		}

		kind := typ.Kind()

		for _, ienc := range encCfg.ifaceEncoder {
			if typ.Implements(ienc.Type) {
				if ifaceIndir(rType) {
					v := &ifaceValEncoder{
						isEmpty: isEmptyFunctions[kind],
						encoder: ienc.Encoder,
						rtype:   rType,
					}
					internalCache[rType] = v
					cache[rType] = v
					continue OuterLoop
				}
				if typ.Kind() == reflect.Ptr {
					elemType := typ.Elem()
					if elemType.Implements(ienc.Type) {
						typesToCreate.push(elemType)
						w := newPointerEncoder(elemType)
						internalCache[rType] = w.encoder
						rebuildMap[rType] = w
					} else {
						v := &pointerIfaceValEncoder{
							encoder: ienc.Encoder,
							rtype:   rType,
						}
						internalCache[rType] = v
						cache[rType] = &directEncoder{v}
					}
					continue OuterLoop
				}
				v := &directIfaceValEncoder{
					isEmpty: isEmptyFunctions[kind],
					encoder: ienc.Encoder,
					rtype:   rType,
				}
				internalCache[rType] = v
				cache[rType] = &directEncoder{v}
				continue OuterLoop
			}
		}

		//// check json.Marshaler interface
		//if typ.Implements(jsonMarshalerType) {
		//	if ifaceIndir(rType) {
		//		v := &jsonMarshalerEncoder{
		//			isEmpty: isEmptyFunctions[kind],
		//			rtype:   rType,
		//		}
		//		internalCache[rType] = v
		//		cache[rType] = v
		//		continue
		//	}
		//	if typ.Kind() == reflect.Ptr {
		//		elemType := typ.Elem()
		//		if elemType.Implements(jsonMarshalerType) {
		//			// treat as a pointer encoder
		//			typesToCreate.push(elemType)
		//			w := newPointerEncoder(elemType)
		//			internalCache[rType] = w.encoder
		//			rebuildMap[rType] = w
		//		} else {
		//			v := pointerJSONMarshalerEncoder(rType)
		//			internalCache[rType] = v
		//			cache[rType] = &directEncoder{v}
		//		}
		//		continue
		//	}
		//	v := &directJSONMarshalerEncoder{
		//		isEmpty: isEmptyFunctions[kind],
		//		rtype:   rType,
		//	}
		//	internalCache[rType] = v
		//	cache[rType] = &directEncoder{v}
		//	continue
		//}
		//
		//// check encoding.TextMarshaler interface
		//if typ.Implements(textMarshalerType) {
		//	if ifaceIndir(rType) {
		//		v := &textMarshalerEncoder{
		//			isEmpty: isEmptyFunctions[kind],
		//			rtype:   rType,
		//		}
		//		internalCache[rType] = v
		//		cache[rType] = v
		//		continue
		//	}
		//	if typ.Kind() == reflect.Ptr {
		//		elemType := typ.Elem()
		//		if elemType.Implements(textMarshalerType) {
		//			// treat as a pointer encoder
		//			typesToCreate.push(elemType)
		//			w := newPointerEncoder(elemType)
		//			internalCache[rType] = w.encoder
		//			rebuildMap[rType] = w
		//		} else {
		//			v := pointerTextMarshalerEncoder(rType)
		//			internalCache[rType] = v
		//			cache[rType] = &directEncoder{v}
		//		}
		//		continue
		//	}
		//	v := &directTextMarshalerEncoder{
		//		isEmpty: isEmptyFunctions[kind],
		//		rtype:   rType,
		//	}
		//	internalCache[rType] = v
		//	cache[rType] = &directEncoder{v}
		//	continue
		//}

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
			w := encCfg.newStructEncoder(typ)
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
				x.encoder.fields.add(fi, v)
			}
			if ifaceIndir(rType) {
				cache[rType] = x.encoder
			} else {
				cache[rType] = &directEncoder{x.encoder}
			}
		}
	}
}
