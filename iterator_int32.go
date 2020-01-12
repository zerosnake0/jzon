package jzon

import (
	"io"
	"math"
	"strconv"
)

const (
	maxUint32Div10NumDigits = 9
	maxUint32Div10          = uint32(math.MaxUint32) / 10
	maxUint32Mod10          = int8(math.MaxUint32 - maxUint32Div10*10)
)

func (it *Iterator) ReadUint32() (uint32, error) {
	c, err := it.nextToken()
	if err != nil {
		return 0, err
	}
	it.head += 1
	return it.readUint32(c)
}

func (it *Iterator) readUint32(c byte) (ret uint32, err error) {
	u := intDigits[c]
	if u == 0 {
		return 0, nil
	}
	if u == invalidDigit {
		return 0, InvalidDigitError{c: c}
	}
	ret = uint32(u)
	numDigit := 1
	for {
		i := it.head
		for ; i < it.tail; i++ {
			digit := intDigits[it.buffer[i]]
			if digit == invalidDigit {
				it.head = i
				return
			}
			if numDigit >= maxUint32Div10NumDigits {
				if ret > maxUint32Div10 ||
					(ret == maxUint32Div10 && digit > maxUint32Mod10) {
					it.head = i
					err = IntOverflowError{}
					return
				}
			}
			ret = (ret << 3) + (ret << 1) + uint32(digit)
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

func (it *Iterator) readInt32(c byte) (int32, error) {
	if c == '-' {
		c, err := it.nextByte()
		if err != nil {
			return 0, err
		}
		it.head += 1
		v, err := it.readUint32(c)
		if err != nil {
			return 0, err
		}
		if v > math.MaxInt32+1 {
			return 0, IntOverflowError{
				typ:   "int32",
				value: "-" + strconv.FormatUint(uint64(v), 10),
			}
		}
		return -int32(v), nil
	} else {
		v, err := it.readUint32(c)
		if err != nil {
			return 0, err
		}
		if v > math.MaxInt32 {
			return 0, IntOverflowError{
				typ:   "int32",
				value: strconv.FormatUint(uint64(v), 10),
			}
		}
		return int32(v), nil
	}
}

func (it *Iterator) ReadInt32() (int32, error) {
	c, err := it.nextToken()
	if err != nil {
		return 0, err
	}
	it.head += 1
	return it.readInt32(c)
}
