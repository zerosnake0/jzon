package jzon

import (
	"io"
	"strconv"
)

func (it *Iterator) ReadFloat32() (float32, error) {
	c, err := it.nextToken()
	if err != nil {
		return 0, err
	}
	it.head += 1
	if c == '-' {
		c, err = it.nextByte()
		if err != nil {
			return 0, err
		}
		it.head += 1
		f, buf, err := it.readPositiveFloat32(c, it.tmpBuffer[:0])
		it.tmpBuffer = buf
		return -f, err
	}
	f, buf, err := it.readPositiveFloat32(c, it.tmpBuffer[:0])
	it.tmpBuffer = buf
	return f, err
}

func (it *Iterator) parseFloat32(buf []byte) (ret float32, err error) {
	f, err := strconv.ParseFloat(localByteToString(buf), 32)
	return float32(f), err
}

func (it *Iterator) readFloat32ExponentPart(buf []byte) (ret float32, _ []byte, err error) {
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
				ret, err = it.parseFloat32(buf)
				return ret, buf, err
			}
		}
		// i == it.tail
		buf = append(buf, it.buffer[it.head:i]...)
		it.head = i
		if err = it.readMore(); err != nil {
			if err == io.EOF {
				ret, err = it.parseFloat32(buf)
				return ret, buf, err
			}
			return 0, buf, err
		}
	}
}

func (it *Iterator) readFloat32FractionPart(buf []byte) (ret float32, _ []byte, err error) {
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
					return it.readFloat32ExponentPart(buf)
				default:
					buf = append(buf, it.buffer[it.head:i]...)
					it.head = i
					ret, err = it.parseFloat32(buf)
					return ret, buf, err
				}
			}
		}
		// i == it.tail
		buf = append(buf, it.buffer[it.head:i]...)
		it.head = i
		if err = it.readMore(); err != nil {
			if err == io.EOF {
				ret, err = it.parseFloat32(buf)
				return ret, buf, err
			}
			return 0, buf, err
		}
	}
}

func (it *Iterator) readPositiveFloat32(c byte, buf []byte) (ret float32, _ []byte, err error) {
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
			return it.readFloat32FractionPart(buf)
		case expInNumber:
			it.head += 1
			buf = append(buf, 'e')
			return it.readFloat32ExponentPart(buf)
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
						return it.readFloat32FractionPart(buf)
					case expInNumber:
						buf = append(buf, it.buffer[it.head:i+1]...)
						it.head = i + 1
						return it.readFloat32ExponentPart(buf)
					default:
						buf = append(buf, it.buffer[it.head:i]...)
						it.head = i
						ret, err = it.parseFloat32(buf)
						return ret, buf, err
					}
				}
			}
			// i == it.tail
			buf = append(buf, it.buffer[it.head:i]...)
			it.head = i
			if err = it.readMore(); err != nil {
				if err == io.EOF {
					ret, err = it.parseFloat32(buf)
					return ret, buf, err
				}
				return 0, buf, err
			}
		}
	}
}
