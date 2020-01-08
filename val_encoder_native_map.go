package jzon

import (
	"encoding"
	"reflect"
	"unsafe"
)

type mapEncoderBuilder struct {
	encoder   *directMapEncoder
	elemRType rtype
}

func newMapEncoder(mapType reflect.Type) *mapEncoderBuilder {
	keyType := mapType.Key()
	var (
		keyEncoder keyEncoder
	)
	keyRType := rtypeOfType(keyType)
	keyKind := keyType.Kind()
	switch {
	case keyKind == reflect.String:
		keyEncoder = keyEncoders[keyKind]
	case keyType.Implements(textMarshalerType):
		if ifaceIndir(keyRType) {
			keyEncoder = textMarshalerKeyEncoder(keyRType)
		} else {
			keyEncoder = directTextMarshalerKeyEncoder(keyRType)
		}
	default:
		if keyEncoder = keyEncoders[keyKind]; keyEncoder == nil {
			return nil
		}
	}
	return &mapEncoderBuilder{
		encoder: &directMapEncoder{
			mapRType: rtypeOfType(mapType),
			// keyRType:   keyRType,
			keyEncoder: keyEncoder,
		},
		elemRType: rtypeOfType(mapType.Elem()),
	}
}

type directMapEncoder struct {
	mapRType rtype

	// keyRType   rtype
	keyEncoder keyEncoder

	elemEncoder ValEncoder
}

func (enc *directMapEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	ptr = *(*unsafe.Pointer)(ptr)
	return maplen(ptr) == 0
}

func (enc *directMapEncoder) Encode(ptr unsafe.Pointer, s *Streamer, _ *EncOpts) {
	if s.Error != nil {
		return
	}
	if ptr == nil {
		s.Null()
		return
	}
	ptr = *(*unsafe.Pointer)(ptr)
	if ptr == nil {
		s.Null()
		return
	}
	s.ObjectStart()
	iter := mapiterinit(enc.mapRType, ptr)
	for i := 0; iter.key != nil; i++ {
		enc.keyEncoder.Encode(iter.key, s)
		if s.Error != nil {
			return
		}
		enc.elemEncoder.Encode(iter.value, s, nil)
		if s.Error != nil {
			return
		}
		mapiternext(iter)
	}
	s.ObjectEnd()
}

// text marshaler
type textMarshalerKeyEncoder rtype

func (enc textMarshalerKeyEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if s.Error != nil {
		return
	}
	rtype := rtype(enc)
	obj := packEFace(rtype, ptr)
	marshaler := obj.(encoding.TextMarshaler)
	b, err := marshaler.MarshalText()
	if err != nil {
		s.Error = err
		return
	}
	s.String(localByteToString(b))
	s.buffer = append(s.buffer, ':')
	s.poped = false
}

type directTextMarshalerKeyEncoder rtype

func (enc directTextMarshalerKeyEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if s.Error != nil {
		return
	}
	rtype := rtype(enc)
	obj := packEFace(rtype, *(*unsafe.Pointer)(ptr))
	marshaler := obj.(encoding.TextMarshaler)
	b, err := marshaler.MarshalText()
	if err != nil {
		s.Error = err
		return
	}
	s.String(localByteToString(b))
	s.buffer = append(s.buffer, ':')
	s.poped = false
}

// key encoders
type stringKeyEncoder struct{}

func (*stringKeyEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	s.Field(*(*string)(ptr))
}

// int encoders
type int8KeyEncoder struct{}

func (*int8KeyEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if s.Error != nil {
		return
	}
	var b [6]byte // `"-128"` max length is 6
	b[0] = '"'
	buf := appendInt8(b[:1], *(*int8)(ptr))
	buf = append(buf, '"')
	s.RawField(buf)
}

type int16KeyEncoder struct{}

func (*int16KeyEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if s.Error != nil {
		return
	}
	var b [8]byte // `"-32768"` max length is 8
	b[0] = '"'
	buf := appendInt16(b[:1], *(*int16)(ptr))
	buf = append(buf, '"')
	s.RawField(buf)
}

type int32KeyEncoder struct{}

func (*int32KeyEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if s.Error != nil {
		return
	}
	var b [13]byte // `"-2147483648"` max length is 13
	b[0] = '"'
	buf := appendInt32(b[:1], *(*int32)(ptr))
	buf = append(buf, '"')
	s.RawField(buf)
}

type int64KeyEncoder struct{}

func (*int64KeyEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if s.Error != nil {
		return
	}
	var b [22]byte // `"-9223372036854775808"` max length is 22
	b[0] = '"'
	buf := appendInt64(b[:1], *(*int64)(ptr))
	buf = append(buf, '"')
	s.RawField(buf)
}

// uint encoders
type uint8KeyEncoder struct{}

func (*uint8KeyEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if s.Error != nil {
		return
	}
	var b [6]byte // `"-128"` max length is 6
	b[0] = '"'
	buf := appendUint8(b[:1], *(*uint8)(ptr))
	buf = append(buf, '"')
	s.RawField(buf)
}

type uint16KeyEncoder struct{}

func (*uint16KeyEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if s.Error != nil {
		return
	}
	var b [8]byte // `"-32768"` max length is 8
	b[0] = '"'
	buf := appendUint16(b[:1], *(*uint16)(ptr))
	buf = append(buf, '"')
	s.RawField(buf)
}

type uint32KeyEncoder struct{}

func (*uint32KeyEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if s.Error != nil {
		return
	}
	var b [13]byte // `"-2147483648"` max length is 13
	b[0] = '"'
	buf := appendUint32(b[:1], *(*uint32)(ptr))
	buf = append(buf, '"')
	s.RawField(buf)
}

type uint64KeyEncoder struct{}

func (*uint64KeyEncoder) Encode(ptr unsafe.Pointer, s *Streamer) {
	if s.Error != nil {
		return
	}
	var b [22]byte // `"-9223372036854775808"` max length is 22
	b[0] = '"'
	buf := appendUint64(b[:1], *(*uint64)(ptr))
	buf = append(buf, '"')
	s.RawField(buf)
}
