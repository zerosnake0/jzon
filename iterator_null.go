package jzon

func (it *Iterator) ReadNull() error {
	c, _, err := it.nextToken()
	if err != nil {
		return err
	}
	if c != 'n' {
		return UnexpectedByteError{got: c, exp: 'n'}
	}
	it.head += 1
	return it.expectBytes("ull")
}
