package jzon

import (
	"io"
	"math"
	"strconv"
)

const (
	maxUint8Div10 = uint8(math.MaxUint8) / 10
	maxUint8Mod10 = int8(math.MaxUint8 - maxUint8Div10*10)
)

func (it *Iterator) ReadUint8() (uint8, error) {
	c, err := it.nextToken()
	if err != nil {
		return 0, err
	}
	it.head += 1
	return it.readUint8(c)
}

func (it *Iterator) readUint8(c byte) (ret uint8, err error) {
	u := intDigits[c]
	if u == 0 {
		return 0, nil
	}
	if u == invalidDigit {
		return 0, InvalidDigitError{c: c}
	}
	ret = uint8(u)
	if it.head == it.tail {
		if err = it.readMore(); err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
	}
	u = intDigits[it.buffer[it.head]] // second digit
	if u == invalidDigit {
		return
	}
	ret = (ret << 3) + (ret << 1) + uint8(u)
	it.head += 1
	if it.head == it.tail {
		if err = it.readMore(); err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
	}
	u = intDigits[it.buffer[it.head]] // third digit
	if u == invalidDigit {
		return
	}
	it.head += 1
	if ret > maxUint8Div10 ||
		(ret == maxUint8Div10 && u > maxUint8Mod10) {
		err = IntOverflowError{}
		return
	}
	ret = (ret << 3) + (ret << 1) + uint8(u)
	return
}

func (it *Iterator) readInt8(c byte) (int8, error) {
	if c == '-' {
		c, err := it.nextByte()
		if err != nil {
			return 0, err
		}
		it.head += 1
		v, err := it.readUint8(c)
		if err != nil {
			return 0, err
		}
		if v > math.MaxInt8+1 {
			return 0, IntOverflowError{
				typ:   "int8",
				value: "-" + strconv.FormatUint(uint64(v), 10),
			}
		}
		return -int8(v), nil
	} else {
		v, err := it.readUint8(c)
		if err != nil {
			return 0, err
		}
		if v > math.MaxInt8 {
			return 0, IntOverflowError{
				typ:   "int8",
				value: strconv.FormatUint(uint64(v), 10),
			}
		}
		return int8(v), nil
	}
}

func (it *Iterator) ReadInt8() (int8, error) {
	c, err := it.nextToken()
	if err != nil {
		return 0, err
	}
	it.head += 1
	return it.readInt8(c)
}
