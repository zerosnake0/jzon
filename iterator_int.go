package jzon

import (
	"strconv"
)

const (
	invalidDigit = -1
)

var (
	intDigits [charNum]int8
)

func init() {
	for i := 0; i < charNum; i++ {
		intDigits[i] = invalidDigit
	}
	for i := '0'; i <= '9'; i++ {
		intDigits[i] = int8(i - '0')
	}
}

func (it *Iterator) ReadInt() (int, error) {
	if strconv.IntSize == 32 {
		i, err := it.ReadInt32()
		return int(i), err
	}
	i, err := it.ReadInt64()
	return int(i), err
}

func (it *Iterator) ReadUint() (uint, error) {
	if strconv.IntSize == 32 {
		u, err := it.ReadUint32()
		return uint(u), err
	}
	u, err := it.ReadUint64()
	return uint(u), err
}
