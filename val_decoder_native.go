package jzon

import (
	"io"
	"unsafe"
)

// bool decoder
type boolDecoder struct{}

func (*boolDecoder) Decode(ptr unsafe.Pointer, it *Iterator, opts *DecOpts) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	switch c {
	case 'n':
		it.head += 1
		return it.expectBytes("ull")
	case 't':
		it.head += 1
		if err := it.expectBytes("rue"); err != nil {
			return err
		}
		*(*bool)(ptr) = true
		return nil
	case 'f':
		it.head += 1
		if err := it.expectBytes("alse"); err != nil {
			return err
		}
		*(*bool)(ptr) = false
		return nil
	case '"':
		quoted := (opts != nil) && opts.Quoted
		if !quoted {
			return UnexpectedByteError{got: c}
		}
		it.head += 1
		c, err := it.nextToken()
		if err != nil {
			return err
		}
		it.head += 1
		switch c {
		case 't':
			if err := it.expectBytes(`rue"`); err != nil {
				return err
			}
			*(*bool)(ptr) = true
			return nil
		case 'f':
			if err := it.expectBytes(`alse"`); err != nil {
				return err
			}
			*(*bool)(ptr) = false
			return nil
		case 'n':
			return it.expectBytes(`ull"`)
		default:
			return UnexpectedByteError{got: c}
		}
	default:
		return UnexpectedByteError{got: c}
	}
}

// string decoder
type stringDecoder struct{}

func (*stringDecoder) Decode(ptr unsafe.Pointer, it *Iterator, opts *DecOpts) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	switch c {
	case '"':
		it.head += 1
		s, err := it.readString()
		if err != nil {
			return err
		}
		quoted := (opts != nil) && opts.Quoted
		if !quoted {
			*(*string)(ptr) = s
			return nil
		}
		l := len(s)
		if l < 2 {
			return BadQuotedStringError(s)
		}
		switch s[0] {
		case 'n':
			if s != "null" {
				return BadQuotedStringError(s)
			}
			return nil
		case '"':
			if s[l-1] != '"' {
				return BadQuotedStringError(s)
			}
			// borrow another iterator
			subIt := it.decoder.NewIterator()
			defer it.decoder.ReturnIterator(subIt)
			subIt.ResetBytes(localStringToBytes(s))
			subStr, err := subIt.ReadString()
			if err != nil {
				return BadQuotedStringError(s)
			}
			// check eof
			_, err = subIt.nextToken()
			if err != io.EOF {
				return BadQuotedStringError(s)
			}
			*(*string)(ptr) = subStr
			return nil
		default:
			return BadQuotedStringError(s)
		}
	case 'n':
		it.head += 1
		return it.expectBytes("ull")
	default:
		return UnexpectedByteError{got: c, exp: '"', exp2: 'n'}
	}
}
