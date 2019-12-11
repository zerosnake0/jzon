package jzon

import (
	"unsafe"
)

// uint8 decoder
type uint8Decoder struct {
}

func (*uint8Decoder) Decode(ptr unsafe.Pointer, it *Iterator) error {
	c, _, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		return it.expectBytes("ull")
	}
	i, err := it.ReadUint8()
	if err != nil {
		return err
	}
	*(*uint8)(ptr) = i
	return nil
}

// uint16 decoder
type uint16Decoder struct {
}

func (*uint16Decoder) Decode(ptr unsafe.Pointer, it *Iterator) error {
	c, _, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		return it.expectBytes("ull")
	}
	i, err := it.ReadUint16()
	if err != nil {
		return err
	}
	*(*uint16)(ptr) = i
	return nil
}

// uint32 decoder
type uint32Decoder struct {
}

func (*uint32Decoder) Decode(ptr unsafe.Pointer, it *Iterator) error {
	c, _, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		return it.expectBytes("ull")
	}
	i, err := it.ReadUint32()
	if err != nil {
		return err
	}
	*(*uint32)(ptr) = i
	return nil
}

// uint64 decoder
type uint64Decoder struct {
}

func (*uint64Decoder) Decode(ptr unsafe.Pointer, it *Iterator) error {
	c, _, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		return it.expectBytes("ull")
	}
	i, err := it.ReadUint64()
	if err != nil {
		return err
	}
	*(*uint64)(ptr) = i
	return nil
}
