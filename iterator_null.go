package jzon

// ReadNull reads a nil
func (it *Iterator) ReadNull() error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	if c != 'n' {
		return UnexpectedByteError{got: c, exp: 'n'}
	}
	it.head++
	return it.expectBytes("ull")
}
