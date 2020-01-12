package jzon

import (
	"unsafe"
)

// float32 decoder
type float32Decoder struct {
}

func (*float32Decoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		return it.expectBytes("ull")
	}
	f, err := it.ReadFloat32()
	if err != nil {
		return err
	}
	*(*float32)(ptr) = f
	return nil
}

// float64 decoder
type float64Decoder struct {
}

func (*float64Decoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		return it.expectBytes("ull")
	}
	f, err := it.ReadFloat64()
	if err != nil {
		return err
	}
	*(*float64)(ptr) = f
	return nil
}
