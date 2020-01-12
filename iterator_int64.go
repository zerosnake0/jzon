package jzon

import (
	"io"
	"math"
	"strconv"
)

const (
	maxUint64Div10NumDigits = 19
	maxUint64Div10          = uint64(math.MaxUint64) / 10
	maxUint64Mod10          = int8(math.MaxUint64 - maxUint64Div10*10)
)

func (it *Iterator) ReadUint64() (uint64, error) {
	c, err := it.nextToken()
	if err != nil {
		return 0, err
	}
	it.head += 1
	return it.readUint64(c)
}

func (it *Iterator) readUint64(c byte) (ret uint64, err error) {
	u := intDigits[c]
	if u == 0 {
		return 0, nil
	}
	if u == invalidDigit {
		return 0, InvalidDigitError{c: c}
	}
	ret = uint64(u)
	numDigit := 1
	// TODO: inline expansion
	for {
		i := it.head
		for ; i < it.tail; i++ {
			digit := intDigits[it.buffer[i]]
			if digit == invalidDigit {
				it.head = i
				return ret, nil
			}
			if numDigit >= maxUint64Div10NumDigits {
				if ret > maxUint64Div10 ||
					(ret == maxUint64Div10 && digit > maxUint64Mod10) {
					it.head = i
					err = IntOverflowError{}
					return
				}
			}
			ret = (ret << 3) + (ret << 1) + uint64(digit)
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

func (it *Iterator) readInt64(c byte) (int64, error) {
	if c == '-' {
		c, err := it.nextByte()
		if err != nil {
			return 0, err
		}
		it.head += 1
		v, err := it.readUint64(c)
		if err != nil {
			return 0, err
		}
		if v > math.MaxInt64+1 {
			return 0, IntOverflowError{
				typ:   "int64",
				value: "-" + strconv.FormatUint(uint64(v), 10),
			}
		}
		return -int64(v), nil
	} else {
		v, err := it.readUint64(c)
		if err != nil {
			return 0, err
		}
		if v > math.MaxInt64 {
			return 0, IntOverflowError{
				typ:   "int64",
				value: strconv.FormatUint(uint64(v), 10),
			}
		}
		return int64(v), nil
	}
}

func (it *Iterator) ReadInt64() (int64, error) {
	c, err := it.nextToken()
	if err != nil {
		return 0, err
	}
	it.head += 1
	return it.readInt64(c)
}
