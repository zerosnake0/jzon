package jzon

// SkipRaw skips and returns the bytes skipped
// slice will not be copied, so make a copy if the return
// value is to be stored somewhere
func (it *Iterator) SkipRaw() ([]byte, error) {
	c, err := it.nextToken()
	if err != nil {
		return nil, err
	}
	oldCapture := it.capture
	it.capture = true
	begin := it.head
	it.head++
	err = skipFunctions[c](it, c)
	it.capture = oldCapture
	if err != nil {
		return nil, err
	}
	return it.buffer[begin:it.head], nil
}

// AppendRaw is like SkipRaw but it is a copy version
func (it *Iterator) AppendRaw(in []byte) ([]byte, error) {
	b, err := it.SkipRaw()
	if err != nil {
		return in, err
	}
	// https://github.com/go101/go101/wiki
	return append(in, b...), nil
}

// ReadRaw reads a raw json object (slice will be copied)
func (it *Iterator) ReadRaw() ([]byte, error) {
	return it.AppendRaw(nil)
}
