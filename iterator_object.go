package jzon

// use after the first `"` is consumed
// will read the object field as well as the colon
func (it *Iterator) readObjectFieldAsSlice() (field []byte, err error) {
	field, err = it.readStringAsSlice()
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
	return
}

// called only when a '"' is consumed
func (it *Iterator) readObjectField() (string, error) {
	field, err := it.readObjectFieldAsSlice()
	if err != nil {
		return "", err
	}
	return string(field), nil
}

func (it *Iterator) skipObjectField() (bool, error) {
	err := skipString(it, '"')
	if err != nil {
		return false, err
	}
	c, err := it.nextToken()
	if err != nil {
		return false, err
	}
	if c != ':' {
		return false, UnexpectedByteError{got: c, exp: ':'}
	}
	it.head += 1
	return true, nil
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
		field, err = it.readObjectField()
		more = err == nil
		return
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
		field, err = it.readObjectField()
		more = err == nil
		return
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
		field, err := it.readObjectField()
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
