package jzon

import (
	"reflect"
	"unsafe"
)

type mapEncoderBuilder struct {
	encoder *directMapEncoder
}

func newMapEncoder(mapType reflect.Type) *mapEncoderBuilder {
	keyType := mapType.Key()
	var (
		keyEncoder ValEncoder
	)
	keyRType := rtypeOfType(keyType)
	if keyEncoder = keyEncoders[keyType.Kind()]; keyEncoder == nil {
		if !keyType.Implements(textMarshalerType) {
			return nil
		}
		keyEncoder = textMarshalerEncoder(keyRType)
	}
	return &mapEncoderBuilder{
		encoder: &directMapEncoder{
			mapRType:   rtypeOfType(mapType),
			keyRType:   keyRType,
			keyEncoder: keyEncoder,
			elemRType:  rtypeOfType(mapType.Elem()),
		},
	}
}

type directMapEncoder struct {
	mapRType rtype

	keyRType   rtype
	keyEncoder ValEncoder

	elemRType   rtype
	elemEncoder ValEncoder
}

func (enc *directMapEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if ptr == nil {
		s.Null()
		return
	}
	s.ObjectStart()
	iter := mapiterinit(enc.mapRType, ptr)
	for i := 0; iter.key != nil; i++ {
		enc.keyEncoder.Encode(iter.key, s)
		enc.elemEncoder.Encode(iter.value, s)
		mapiternext(iter)
	}
	s.ObjectEnd()
}

// key encoders
type stringKeyEncoder struct{}

func (enc *stringKeyEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	s.Field(*(*string)(ptr))
}
