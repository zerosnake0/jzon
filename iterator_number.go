package jzon

// it.head has already be forwarded before calling this
func (it *Iterator) readNumberAsString(c byte) (n string, err error) {
	// start capture
	oldCapture := it.capture
	it.capture = true

	// save current location
	// this should be a safe operation because the c is still in the buffer
	begin := it.head - 1
	err = skipNumber(it, c)
	if err == nil {
		n = string(it.buffer[begin:it.head])
	}
	// end capture
	it.capture = oldCapture
	return
}

func (it *Iterator) ReadNumber() (n Number, err error) {
	c, err := it.nextToken()
	if err != nil {
		return
	}
	if valueTypeMap[c] != NumberValue {
		err = UnexpectedByteError{got: c}
		return
	}
	it.head += 1
	s, err := it.readNumberAsString(c)
	return Number(s), err
}
