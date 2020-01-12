package jzon

import (
	"io"
	"math"
	"strconv"
)

const (
	maxUint16Div10NumDigits = 4
	maxUint16Div10          = uint16(math.MaxUint16) / 10
	maxUint16Mod10          = int8(math.MaxUint16 - maxUint16Div10*10)
)

func (it *Iterator) ReadUint16() (uint16, error) {
	c, err := it.nextToken()
	if err != nil {
		return 0, err
	}
	it.head += 1
	return it.readUint16(c)
}

func (it *Iterator) readUint16(c byte) (ret uint16, err error) {
	u := intDigits[c]
	if u == 0 {
		return 0, nil
	}
	if u == invalidDigit {
		return 0, InvalidDigitError{c: c}
	}
	ret = uint16(u)
	numDigit := 1
	for {
		i := it.head
		for ; i < it.tail; i++ {
			digit := intDigits[it.buffer[i]]
			if digit == invalidDigit {
				it.head = i
				return ret, nil
			}
			if numDigit >= maxUint16Div10NumDigits {
				if ret > maxUint16Div10 ||
					(ret == maxUint16Div10 && digit > maxUint16Mod10) {
					it.head = i
					err = IntOverflowError{}
					return
				}
			}
			ret = (ret << 3) + (ret << 1) + uint16(digit)
			numDigit++
		}
		it.head = i
		if err = it.readMore(); err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
	}
}

func (it *Iterator) readInt16(c byte) (int16, error) {
	if c == '-' {
		c, err := it.nextByte()
		if err != nil {
			return 0, err
		}
		it.head += 1
		v, err := it.readUint16(c)
		if err != nil {
			return 0, err
		}
		if v > math.MaxInt16+1 {
			return 0, IntOverflowError{
				typ:   "int16",
				value: "-" + strconv.FormatUint(uint64(v), 10),
			}
		}
		return -int16(v), nil
	} else {
		v, err := it.readUint16(c)
		if err != nil {
			return 0, err
		}
		if v > math.MaxInt16 {
			return 0, IntOverflowError{
				typ:   "int16",
				value: strconv.FormatUint(uint64(v), 10),
			}
		}
		return int16(v), nil
	}
}

func (it *Iterator) ReadInt16() (int16, error) {
	c, err := it.nextToken()
	if err != nil {
		return 0, err
	}
	it.head += 1
	return it.readInt16(c)
}
