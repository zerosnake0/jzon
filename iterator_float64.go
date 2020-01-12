package jzon

import (
	"io"
	"strconv"
)

func (it *Iterator) readFloat64(c byte) (float64, error) {
	if c == '-' {
		c, err := it.nextByte()
		if err != nil {
			return 0, err
		}
		it.head += 1
		f, buf, err := it.readPositiveFloat64(c, it.tmpBuffer[:0])
		it.tmpBuffer = buf
		return -f, err
	}
	f, buf, err := it.readPositiveFloat64(c, it.tmpBuffer[:0])
	it.tmpBuffer = buf
	return f, err
}

func (it *Iterator) ReadFloat64() (float64, error) {
	c, err := it.nextToken()
	if err != nil {
		return 0, err
	}
	it.head += 1
	return it.readFloat64(c)
}

func (it *Iterator) parseFloat64(buf []byte) (ret float64, err error) {
	return strconv.ParseFloat(localByteToString(buf), 64)
}

func (it *Iterator) readFloat64ExponentPart(buf []byte) (ret float64, _ []byte, err error) {
	c, err := it.nextByte()
	if err != nil {
		return 0, buf, err
	}
	it.head += 1
	if c == '+' || c == '-' {
		buf = append(buf, c)
		c, err = it.nextByte()
		if err != nil {
			return 0, buf, err
		}
		it.head += 1
	}
	if intDigits[c] == invalidDigit {
		return 0, buf, InvalidFloatError{c: c}
	}
	buf = append(buf, c)
	for {
		i := it.head
		for ; i < it.tail; i++ {
			c = it.buffer[i]
			digit := intDigits[c]
			if digit == invalidDigit {
				buf = append(buf, it.buffer[it.head:i]...)
				it.head = i
				f, err := it.parseFloat64(buf)
				return f, buf, err
			}
		}
		// i == it.tail
		buf = append(buf, it.buffer[it.head:i]...)
		it.head = i
		if err = it.readMore(); err != nil {
			if err == io.EOF {
				f, err := it.parseFloat64(buf)
				return f, buf, err
			}
			return 0, buf, err
		}
	}
}

func (it *Iterator) readFloat64FractionPart(buf []byte) (ret float64, _ []byte, err error) {
	c, err := it.nextByte()
	if err != nil {
		return 0, buf, err
	}
	if intDigits[c] == invalidDigit {
		return 0, buf, InvalidFloatError{c: c}
	}
	it.head += 1
	buf = append(buf, c)
	for {
		i := it.head
		for ; i < it.tail; i++ {
			c = it.buffer[i]
			digit := floatDigits[c]
			if digit < 0 {
				switch digit {
				case expInNumber:
					buf = append(buf, it.buffer[it.head:i+1]...)
					it.head = i + 1
					return it.readFloat64ExponentPart(buf)
				default:
					buf = append(buf, it.buffer[it.head:i]...)
					it.head = i
					f, err := it.parseFloat64(buf)
					return f, buf, err
				}
			}
		}
		// i == it.tail
		buf = append(buf, it.buffer[it.head:i]...)
		it.head = i
		if err = it.readMore(); err != nil {
			if err == io.EOF {
				f, err := it.parseFloat64(buf)
				return f, buf, err
			}
			return 0, buf, err
		}
	}
}

func (it *Iterator) readPositiveFloat64(c byte, buf []byte) (ret float64, _ []byte, err error) {
	u := intDigits[c]
	if u == invalidDigit {
		return 0, buf, InvalidFloatError{c: c}
	}

	buf = append(buf, c)
	if u == 0 {
		if it.head == it.tail {
			if err = it.readMore(); err != nil {
				if err == io.EOF {
					return 0, buf, nil
				}
				return 0, buf, err
			}
		}
		switch floatDigits[it.buffer[it.head]] {
		case dotInNumber:
			it.head += 1
			buf = append(buf, '.')
			return it.readFloat64FractionPart(buf)
		case expInNumber:
			it.head += 1
			buf = append(buf, 'e')
			return it.readFloat64ExponentPart(buf)
		default:
			return 0, buf, nil
		}
	} else {
		for {
			i := it.head
			for ; i < it.tail; i++ {
				c = it.buffer[i]
				digit := floatDigits[c]
				if digit < 0 {
					switch digit {
					case dotInNumber:
						buf = append(buf, it.buffer[it.head:i+1]...)
						it.head = i + 1
						return it.readFloat64FractionPart(buf)
					case expInNumber:
						buf = append(buf, it.buffer[it.head:i+1]...)
						it.head = i + 1
						return it.readFloat64ExponentPart(buf)
					default:
						buf = append(buf, it.buffer[it.head:i]...)
						it.head = i
						f, err := it.parseFloat64(buf)
						return f, buf, err
					}
				}
			}
			// i == it.tail
			buf = append(buf, it.buffer[it.head:i]...)
			it.head = i
			if err = it.readMore(); err != nil {
				if err == io.EOF {
					f, err := it.parseFloat64(buf)
					return f, buf, err
				}
				return 0, buf, err
			}
		}
	}
}
