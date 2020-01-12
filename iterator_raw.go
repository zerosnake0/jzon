package jzon

// No copy version
func (it *Iterator) SkipRaw() ([]byte, error) {
	c, err := it.nextToken()
	if err != nil {
		return nil, err
	}
	oldCapture := it.capture
	it.capture = true
	begin := it.head
	it.head += 1
	err = skipFunctions[c](it, c)
	it.capture = oldCapture
	if err != nil {
		return nil, err
	}
	return it.buffer[begin:it.head], nil
}

// copy version
func (it *Iterator) AppendRaw(in []byte) ([]byte, error) {
	b, err := it.SkipRaw()
	if err != nil {
		return in, err
	}
	// https://github.com/go101/go101/wiki
	return append(in, b...), nil
}

// copy version
func (it *Iterator) ReadRaw() ([]byte, error) {
	return it.AppendRaw(nil)
}
