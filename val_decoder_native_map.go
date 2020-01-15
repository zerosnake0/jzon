package jzon

import (
	"reflect"
	"unsafe"
)

type mapDecoderBuilder struct {
	decoder     *mapDecoder
	valPtrRType rtype
}

func newMapDecoder(mapType reflect.Type) *mapDecoderBuilder {
	// Compatible with standard lib
	// Map key must either have string kind, have an integer kind,
	// or be an encoding.TextUnmarshaler.
	keyType := mapType.Key()
	var (
		keyDecoder ValDecoder
	)
	keyPtrType := reflect.PtrTo(keyType)
	keyKind := keyType.Kind()
	// the string type is specially treated in order to be
	// compatible with the standard lib
	switch {
	case keyKind == reflect.String:
		keyDecoder = keyDecoders[keyKind]
	case keyPtrType.Implements(textUnmarshalerType):
		keyDecoder = textUnmarshalerDecoder(rtypeOfType(keyPtrType))
	default:
		if keyDecoder = keyDecoders[keyType.Kind()]; keyDecoder == nil {
			return nil
		}
	}
	return &mapDecoderBuilder{
		decoder: &mapDecoder{
			rtype:    rtypeOfType(mapType),
			keyRType: rtypeOfType(keyType),
			keyDec:   keyDecoder,
			valRType: rtypeOfType(mapType.Elem()),
		},
		valPtrRType: rtypeOfType(reflect.PtrTo(mapType.Elem())),
	}
}

func (builder *mapDecoderBuilder) build(cache decoderCache) {
	builder.decoder.valDec = cache[builder.valPtrRType]
}

type mapDecoder struct {
	rtype rtype

	keyRType rtype
	keyDec   ValDecoder

	valRType rtype
	valDec   ValDecoder
}

func (dec *mapDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		if err = it.expectBytes("ull"); err != nil {
			return err
		}
		*(*unsafe.Pointer)(ptr) = nil
		return nil
	}
	if c != '{' {
		return UnexpectedByteError{got: c, exp: '{', exp2: 'n'}
	}
	it.head += 1
	c, err = it.nextToken()
	if err != nil {
		return err
	}
	if c == '}' {
		it.head += 1
		if *(*unsafe.Pointer)(ptr) == nil {
			typedmemmove(dec.rtype, ptr, unsafeMakeMap(dec.rtype, 0))
		}
		return nil
	}
	if *(*unsafe.Pointer)(ptr) == nil {
		typedmemmove(dec.rtype, ptr, unsafeMakeMap(dec.rtype, 0))
	}
	for {
		key := unsafe_New(dec.keyRType)
		opt := DecOpts{
			MapKey: true,
		}
		if err = dec.keyDec.Decode(key, it, opt.noescape()); err != nil {
			return err
		}
		c, err = it.nextToken()
		if err != nil {
			return err
		}
		if c != ':' {
			return UnexpectedByteError{got: c, exp: ':'}
		}
		it.head += 1
		val := unsafe_New(dec.valRType)
		if err = dec.valDec.Decode(val, it, nil); err != nil {
			return err
		}
		mapassign(dec.rtype, *(*unsafe.Pointer)(ptr), key, val)
		c, err = it.nextToken()
		if err != nil {
			return err
		}
		switch c {
		case '}':
			it.head += 1
			return nil
		case ',':
			it.head += 1
		default:
			return UnexpectedByteError{got: c, exp: '}', exp2: ','}
		}
	}
}

/*
// key decoders
type stringKeyDecoder struct {
}

func (*stringKeyDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	s, err := it.ReadString()
	if err != nil {
		return err
	}
	*(*string)(ptr) = s
	return nil
}

// int key decoders
type int8KeyDecoder struct {
}

func (*int8KeyDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	if err := it.expectQuote(); err != nil {
		return err
	}
	i, err := it.ReadInt8()
	if err != nil {
		return err
	}
	if err := it.expectQuote(); err != nil {
		return err
	}
	*(*int8)(ptr) = i
	return nil
}

type int16KeyDecoder struct {
}

func (*int16KeyDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	if err := it.expectQuote(); err != nil {
		return err
	}
	i, err := it.ReadInt16()
	if err != nil {
		return err
	}
	if err := it.expectQuote(); err != nil {
		return err
	}
	*(*int16)(ptr) = i
	return nil
}

type int32KeyDecoder struct {
}

func (*int32KeyDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	if err := it.expectQuote(); err != nil {
		return err
	}
	i, err := it.ReadInt32()
	if err != nil {
		return err
	}
	if err := it.expectQuote(); err != nil {
		return err
	}
	*(*int32)(ptr) = i
	return nil
}

type int64KeyDecoder struct {
}

func (*int64KeyDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	if err := it.expectQuote(); err != nil {
		return err
	}
	i, err := it.ReadInt64()
	if err != nil {
		return err
	}
	if err := it.expectQuote(); err != nil {
		return err
	}
	*(*int64)(ptr) = i
	return nil
}

// uint key decoders
type uint8KeyDecoder struct {
}

func (*uint8KeyDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	if err := it.expectQuote(); err != nil {
		return err
	}
	i, err := it.ReadUint8()
	if err != nil {
		return err
	}
	if err := it.expectQuote(); err != nil {
		return err
	}
	*(*uint8)(ptr) = i
	return nil
}

type uint16KeyDecoder struct {
}

func (*uint16KeyDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	if err := it.expectQuote(); err != nil {
		return err
	}
	i, err := it.ReadUint16()
	if err != nil {
		return err
	}
	if err := it.expectQuote(); err != nil {
		return err
	}
	*(*uint16)(ptr) = i
	return nil
}

type uint32KeyDecoder struct {
}

func (*uint32KeyDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	if err := it.expectQuote(); err != nil {
		return err
	}
	i, err := it.ReadUint32()
	if err != nil {
		return err
	}
	if err := it.expectQuote(); err != nil {
		return err
	}
	*(*uint32)(ptr) = i
	return nil
}

type uint64KeyDecoder struct {
}

func (*uint64KeyDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	if err := it.expectQuote(); err != nil {
		return err
	}
	i, err := it.ReadUint64()
	if err != nil {
		return err
	}
	if err := it.expectQuote(); err != nil {
		return err
	}
	*(*uint64)(ptr) = i
	return nil
}
*/
