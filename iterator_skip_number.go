package jzon

import (
	"io"
)

func (it *Iterator) skipExponentPart() (err error) {
	c, err := it.nextByte()
	if err != nil {
		return err
	}
	it.head += 1
	if c == '+' || c == '-' {
		c, err = it.nextByte()
		if err != nil {
			return err
		}
		it.head += 1
	}
	if intDigits[c] == invalidDigit {
		return InvalidDigitError{c: c}
	}
	for {
		i := it.head
		for ; i < it.tail; i++ {
			c = it.buffer[i]
			digit := intDigits[c]
			if digit == invalidDigit {
				it.head = i
				return nil
			}
		}
		// i == it.tail
		it.head = i
		if err = it.readMore(); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func (it *Iterator) skipFractionPart() (err error) {
	c, err := it.nextByte()
	if err != nil {
		return err
	}
	if intDigits[c] == invalidDigit {
		return InvalidDigitError{c: c}
	}
	it.head += 1
	for {
		i := it.head
		for ; i < it.tail; i++ {
			c = it.buffer[i]
			digit := floatDigits[c]
			if digit < 0 {
				switch digit {
				case expInNumber:
					it.head = i + 1
					return it.skipExponentPart()
				default:
					it.head = i
					return nil
				}
			}
		}
		// i == it.tail
		it.head = i
		if err = it.readMore(); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

// make sure that c is '-' or '0'~'9' before calling this method
// (the call of nextToken before this method return NumberValue)
func skipNumber(it *Iterator, c byte) (err error) {
	if c == '-' {
		c, err = it.nextByte()
		if err != nil {
			return
		}
		if intDigits[c] == invalidDigit {
			return InvalidDigitError{c: c}
		}
		it.head += 1
	}
	// positive
	// here the c can only be '0'~'9'
	if c == '0' {
		if it.head == it.tail {
			if err = it.readMore(); err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
		}
		switch floatDigits[it.buffer[it.head]] {
		case dotInNumber:
			it.head += 1
			return it.skipFractionPart()
		case expInNumber:
			it.head += 1
			return it.skipExponentPart()
		default:
			return nil
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
						it.head = i + 1
						return it.skipFractionPart()
					case expInNumber:
						it.head = i + 1
						return it.skipExponentPart()
					default:
						it.head = i
						return nil
					}
				}
			}
			// i == it.tail
			it.head = i
			if err = it.readMore(); err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
		}
	}
}

func (it *Iterator) SkipNumber() error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	if valueTypeMap[c] != NumberValue {
		return UnexpectedByteError{got: c}
	}
	it.head += 1
	return skipNumber(it, c)
}
