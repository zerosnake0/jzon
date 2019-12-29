package jzon

import (
	"unsafe"
)

// int8 decoder
type int8Decoder struct {
}

func (*int8Decoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	c, _, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		return it.expectBytes("ull")
	}
	i, err := it.ReadInt8()
	if err != nil {
		return err
	}
	*(*int8)(ptr) = i
	return nil
}

// int16 decoder
type int16Decoder struct {
}

func (*int16Decoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	c, _, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		return it.expectBytes("ull")
	}
	i, err := it.ReadInt16()
	if err != nil {
		return err
	}
	*(*int16)(ptr) = i
	return nil
}

// int32 decoder
type int32Decoder struct {
}

func (*int32Decoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	c, _, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		return it.expectBytes("ull")
	}
	i, err := it.ReadInt32()
	if err != nil {
		return err
	}
	*(*int32)(ptr) = i
	return nil
}

// int64 decoder
type int64Decoder struct {
}

func (*int64Decoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	c, _, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		return it.expectBytes("ull")
	}
	i, err := it.ReadInt64()
	if err != nil {
		return err
	}
	*(*int64)(ptr) = i
	return nil
}
