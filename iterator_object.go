package jzon

// use after the first `"` is consumed
// will read the object field as well as the colon
// more must be true when err is nil
func (it *Iterator) readObjectFieldAsSlice(buf []byte) (
	more bool, field []byte, err error) {
	field, err = it.readStringAsSlice(buf)
	if err != nil {
		return
	}
	c, err := it.nextToken()
	if err != nil {
		return
	}
	if c != ':' {
		err = UnexpectedByteError{got: c, exp: ':'}
		return
	}
	it.head += 1
	more = true
	return
}

func (it *Iterator) readObjectField() (bool, string, error) {
	more, field, err := it.readObjectFieldAsSlice(it.tmpBuffer[:0])
	it.tmpBuffer = field
	if err != nil {
		return false, "", err
	}
	return more, string(field), nil
}

func (it *Iterator) skipObjectField() (bool, error) {
	more, field, err := it.readObjectFieldAsSlice(it.tmpBuffer[:0])
	it.tmpBuffer = field
	return more, err
}

func (it *Iterator) ReadObjectBegin() (more bool, field string, err error) {
	c, err := it.nextToken()
	if err != nil {
		return
	}
	if c != '{' {
		err = UnexpectedByteError{got: c, exp: '{'}
		return
	}
	it.head += 1
	c, err = it.nextToken()
	if err != nil {
		return
	}
	switch c {
	case '}':
		// no more items
		it.head += 1
		return
	case '"':
		it.head += 1
		return it.readObjectField()
	default:
		err = UnexpectedByteError{got: c, exp: '}', exp2: '"'}
		return
	}
}

func (it *Iterator) ReadObjectMore() (more bool, field string, err error) {
	c, err := it.nextToken()
	if err != nil {
		return
	}
	switch c {
	case '}':
		it.head += 1
		return
	case ',':
		it.head += 1
		c, err = it.nextToken()
		if err != nil {
			return
		}
		if c != '"' {
			err = UnexpectedByteError{got: c, exp: '"'}
			return
		}
		it.head += 1
		return it.readObjectField()
	default:
		err = UnexpectedByteError{got: c, exp: '}', exp2: ','}
		return
	}
}

func (it *Iterator) ReadObjectCB(cb func(it *Iterator, field string) error) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	if c != '{' {
		return UnexpectedByteError{got: c, exp: '{'}
	}
	it.head += 1
	c, err = it.nextToken()
	if err != nil {
		return err
	}
	if c == '}' {
		it.head += 1
		return nil
	}
	if c != '"' {
		return UnexpectedByteError{got: c, exp: '"', exp2: '}'}
	}
	it.head += 1
	for {
		_, field, err := it.readObjectField()
		if err != nil {
			return err
		}
		if err := cb(it, field); err != nil {
			return err
		}
		c, err = it.nextToken()
		if err != nil {
			return err
		}
		switch c {
		case '}':
			it.head += 1
			return nil
		case ',':
			it.head += 1
			c, err = it.nextToken()
			if err != nil {
				return err
			}
			if c != '"' {
				return UnexpectedByteError{got: c, exp: '"'}
			}
			it.head += 1
		default:
			return UnexpectedByteError{got: c, exp: '}', exp2: ','}
		}
	}
}
