package jzon

import (
	"unsafe"
)

// uint8 decoder
type uint8Decoder struct {
}

func (*uint8Decoder) Decode(ptr unsafe.Pointer, it *Iterator, opts *DecOpts) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		return it.expectBytes("ull")
	}
	quoted := (opts != nil) && (opts.Quoted || opts.MapKey)
	if quoted {
		if c != '"' {
			return UnexpectedByteError{got: c, exp: '"'}
		}
		it.head += 1
		c, err = it.nextByte()
		if err != nil {
			return err
		}
	}
	it.head += 1
	i, err := it.readUint8(c)
	if err != nil {
		return err
	}
	if quoted {
		c, err = it.nextByte()
		if err != nil {
			return err
		}
		if c != '"' {
			return UnexpectedByteError{got: c, exp: '"'}
		}
		it.head += 1
	}
	*(*uint8)(ptr) = i
	return nil
}

// uint16 decoder
type uint16Decoder struct {
}

func (*uint16Decoder) Decode(ptr unsafe.Pointer, it *Iterator, opts *DecOpts) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		return it.expectBytes("ull")
	}
	quoted := (opts != nil) && (opts.Quoted || opts.MapKey)
	if quoted {
		if c != '"' {
			return UnexpectedByteError{got: c, exp: '"'}
		}
		it.head += 1
		c, err = it.nextByte()
		if err != nil {
			return err
		}
	}
	it.head += 1
	i, err := it.readUint16(c)
	if err != nil {
		return err
	}
	if quoted {
		c, err = it.nextByte()
		if err != nil {
			return err
		}
		if c != '"' {
			return UnexpectedByteError{got: c, exp: '"'}
		}
		it.head += 1
	}
	*(*uint16)(ptr) = i
	return nil
}

// uint32 decoder
type uint32Decoder struct {
}

func (*uint32Decoder) Decode(ptr unsafe.Pointer, it *Iterator, opts *DecOpts) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		return it.expectBytes("ull")
	}
	quoted := (opts != nil) && (opts.Quoted || opts.MapKey)
	if quoted {
		if c != '"' {
			return UnexpectedByteError{got: c, exp: '"'}
		}
		it.head += 1
		c, err = it.nextByte()
		if err != nil {
			return err
		}
	}
	it.head += 1
	i, err := it.readUint32(c)
	if err != nil {
		return err
	}
	if quoted {
		c, err = it.nextByte()
		if err != nil {
			return err
		}
		if c != '"' {
			return UnexpectedByteError{got: c, exp: '"'}
		}
		it.head += 1
	}
	*(*uint32)(ptr) = i
	return nil
}

// uint64 decoder
type uint64Decoder struct {
}

func (*uint64Decoder) Decode(ptr unsafe.Pointer, it *Iterator, opts *DecOpts) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == 'n' {
		it.head += 1
		return it.expectBytes("ull")
	}
	quoted := (opts != nil) && (opts.Quoted || opts.MapKey)
	if quoted {
		if c != '"' {
			return UnexpectedByteError{got: c, exp: '"'}
		}
		it.head += 1
		c, err = it.nextByte()
		if err != nil {
			return err
		}
	}
	it.head += 1
	i, err := it.readUint64(c)
	if err != nil {
		return err
	}
	if quoted {
		c, err = it.nextByte()
		if err != nil {
			return err
		}
		if c != '"' {
			return UnexpectedByteError{got: c, exp: '"'}
		}
		it.head += 1
	}
	*(*uint64)(ptr) = i
	return nil
}
