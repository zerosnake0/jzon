package jzon

import (
	"strconv"
)

const (
	invalidDigit = -1
)

var (
	readInt   func(it *Iterator) (int, error)
	readUint  func(it *Iterator) (uint, error)
	intDigits [charNum]int8
)

func init() {
	if strconv.IntSize == 32 {
		readInt = func(it *Iterator) (int, error) {
			i, err := it.ReadInt32()
			return int(i), err
		}
		readUint = func(it *Iterator) (uint, error) {
			u, err := it.ReadUint32()
			return uint(u), err
		}
	} else {
		readInt = func(it *Iterator) (int, error) {
			i, err := it.ReadInt64()
			return int(i), err
		}
		readUint = func(it *Iterator) (uint, error) {
			u, err := it.ReadUint64()
			return uint(u), err
		}
	}
	for i := 0; i < charNum; i++ {
		intDigits[i] = invalidDigit
	}
	for i := '0'; i <= '9'; i++ {
		intDigits[i] = int8(i - '0')
	}
}

func (it *Iterator) ReadInt() (int, error) {
	return readInt(it)
}

func (it *Iterator) ReadUint() (uint, error) {
	return readUint(it)
}
