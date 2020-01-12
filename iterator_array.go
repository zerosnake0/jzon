package jzon

/*
 * var (
 *     more bool
 *     err error
 * )
 * for more, err = it.ReadArray();
 *     more;
 *     more, err = it.ReadArrayMore() {
 * }
 * if err != nil {
 *     // error handling
 * }
 */
func (it *Iterator) ReadArrayBegin() (ret bool, err error) {
	c, err := it.nextToken()
	if err != nil {
		return false, err
	}
	if c != '[' {
		return false, UnexpectedByteError{got: c, exp: '['}
	}
	it.head += 1
	c, err = it.nextToken()
	if err != nil {
		return false, err
	}
	if c == ']' {
		it.head += 1
		return false, nil
	}
	return true, nil
}

func (it *Iterator) ReadArrayMore() (ret bool, err error) {
	c, err := it.nextToken()
	if err != nil {
		return false, err
	}
	switch c {
	case ',':
		it.head += 1
		return true, nil
	case ']':
		it.head += 1
		return false, nil
	default:
		return false, UnexpectedByteError{got: c, exp: ',', exp2: ']'}
	}
}

func (it *Iterator) ReadArrayCB(cb func(*Iterator) error) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	if c != '[' {
		return UnexpectedByteError{got: c, exp: '['}
	}
	it.head += 1
	c, err = it.nextToken()
	if err != nil {
		return err
	}
	if c == ']' {
		it.head += 1
		return nil
	}
	for {
		if err := cb(it); err != nil {
			return err
		}
		c, err = it.nextToken()
		if err != nil {
			return err
		}
		switch c {
		case ',':
			it.head += 1
		case ']':
			it.head += 1
			return nil
		default:
			return UnexpectedByteError{got: c, exp: ',', exp2: ']'}
		}
	}
}
