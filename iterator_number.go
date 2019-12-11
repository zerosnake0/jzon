package jzon

// do not step forward it.head before calling this
func (it *Iterator) readNumberAsString(c byte) (n string, err error) {
	// start capture
	oldCapture := it.capture
	it.capture = true

	begin := it.head // save current location
	it.head += 1
	err = skipNumber(it, c)
	if err == nil {
		n = string(it.buffer[begin:it.head])
	}
	// end capture
	it.capture = oldCapture
	return
}

func (it *Iterator) ReadNumber() (n Number, err error) {
	c, vt, err := it.nextToken()
	if err != nil {
		return
	}
	if vt != NumberValue {
		err = UnexpectedByteError{got: c}
		return
	}
	s, err := it.readNumberAsString(c)
	return Number(s), err
}
