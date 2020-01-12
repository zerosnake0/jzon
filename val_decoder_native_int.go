package jzon

import (
	"unsafe"
)

// int8 decoder
type int8Decoder struct {
}

func (*int8Decoder) Decode(ptr unsafe.Pointer, it *Iterator, opts *DecOpts) error {
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
	i, err := it.readInt8(c)
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
	*(*int8)(ptr) = i
	return nil
}

// int16 decoder
type int16Decoder struct {
}

func (*int16Decoder) Decode(ptr unsafe.Pointer, it *Iterator, opts *DecOpts) error {
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
	i, err := it.readInt16(c)
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
	*(*int16)(ptr) = i
	return nil
}

// int32 decoder
type int32Decoder struct {
}

func (*int32Decoder) Decode(ptr unsafe.Pointer, it *Iterator, opts *DecOpts) error {
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
	i, err := it.readInt32(c)
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
	*(*int32)(ptr) = i
	return nil
}

// int64 decoder
type int64Decoder struct {
}

func (*int64Decoder) Decode(ptr unsafe.Pointer, it *Iterator, opts *DecOpts) error {
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
	i, err := it.readInt64(c)
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
	*(*int64)(ptr) = i
	return nil
}
