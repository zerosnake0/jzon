package jzon

// ReadArrayBegin starts to read an array
//
// var (
//     more bool
//     err error
// )
// for more, err = it.ReadArray();
//     more;
//     more, err = it.ReadArrayMore() {
// }
// if err != nil {
//     // error handling
// }
func (it *Iterator) ReadArrayBegin() (ret bool, err error) {
	c, err := it.nextToken()
	if err != nil {
		return false, err
	}
	if c != '[' {
		return false, UnexpectedByteError{got: c, exp: '['}
	}
	it.head++
	c, err = it.nextToken()
	if err != nil {
		return false, err
	}
	if c == ']' {
		it.head++
		return false, nil
	}
	return true, nil
}

// ReadArrayMore tells if there is more item to read in the array
func (it *Iterator) ReadArrayMore() (ret bool, err error) {
	c, err := it.nextToken()
	if err != nil {
		return false, err
	}
	switch c {
	case ',':
		it.head++
		return true, nil
	case ']':
		it.head++
		return false, nil
	default:
		return false, UnexpectedByteError{got: c, exp: ',', exp2: ']'}
	}
}

// ReadArrayCB reads the array with a callback
// The caller should make sure that the callback is correct
func (it *Iterator) ReadArrayCB(cb func(*Iterator) error) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	if c != '[' {
		return UnexpectedByteError{got: c, exp: '['}
	}
	it.head++
	c, err = it.nextToken()
	if err != nil {
		return err
	}
	if c == ']' {
		it.head++
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
			it.head++
		case ']':
			it.head++
			return nil
		default:
			return UnexpectedByteError{got: c, exp: ',', exp2: ']'}
		}
	}
}
