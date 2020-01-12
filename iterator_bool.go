package jzon

func (it *Iterator) ReadBool() (bool, error) {
	c, err := it.nextToken()
	if err != nil {
		return false, err
	}
	switch c {
	case 't':
		it.head += 1
		return true, it.expectBytes("rue")
	case 'f':
		it.head += 1
		return false, it.expectBytes("alse")
	default:
		return false, UnexpectedByteError{got: c, exp: 't', exp2: 'f'}
	}
}
