package jzon

func (it *Iterator) skipU4() error {
	remain := 4
	for {
		i := it.head
		for ; i < it.tail; i++ {
			c := it.buffer[i]
			u4v := hexValue[c]
			if u4v == invalidHex {
				return InvalidUnicodeCharError{c: c}
			}
			if remain == 1 {
				it.head = i + 1
				return nil
			}
			remain--
		}
		it.head = i
		if err := it.readMore(); err != nil {
			return err
		}
	}
}

func (it *Iterator) skipEscapedChar() error {
	c, err := it.nextByte()
	if err != nil {
		return err
	}
	escaped := escapeMap[c]
	if escaped != noEscape {
		it.head += 1
		return nil
	}
	if c != 'u' {
		return InvalidEscapeCharError{c: c}
	}
	it.head += 1
	return it.skipU4()
}

// internal, call only after a '"' is consumed
func skipString(it *Iterator, _ byte) error {
	for {
		i := it.head
		for i < it.tail {
			c := it.buffer[i]
			if c == '"' {
				it.head = i + 1
				return nil
			} else if c == '\\' {
				it.head = i + 1
				err := it.skipEscapedChar()
				if err != nil {
					return err
				}
				i = it.head
			} else if c < ' ' { // json.org
				return InvalidStringCharError{c: c}
			} else {
				i++
			}
		}
		// i == it.tail
		it.head = i
		if err := it.readMore(); err != nil {
			return err
		}
	}
}

func (it *Iterator) SkipString() error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	if c != '"' {
		return UnexpectedByteError{got: c, exp: '"'}
	}
	it.head += 1
	return skipString(it, c)
}
