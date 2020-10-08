package jzon

func skipObjectWithStack(it *Iterator, _ byte) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == '}' {
		it.head++
		return nil
	}
	for {
		if c != '"' {
			return UnexpectedByteError{got: c, exp: '"'}
		}
		it.head++
		if err = it.skipObjectField(); err != nil {
			return err
		}
		c, err = it.nextToken()
		if err != nil {
			return err
		}
		it.head++
		switch c {
		case '{':
			s := stackPool.Get().(*stack).initObject()
			err := skipWithStack(it, stackElementObjectBegin, s)
			stackPool.Put(s)
			return err
		case '[':
			s := stackPool.Get().(*stack).initObject()
			err := skipWithStack(it, stackElementArrayBegin, s)
			stackPool.Put(s)
			return err
		}
		if err = skipFunctions[c](it, c); err != nil {
			return err
		}
		c, err = it.nextToken()
		if err != nil {
			return err
		}
		it.head++
		if c == '}' {
			return nil
		}
		if c != ',' {
			return UnexpectedByteError{got: c, exp: '}', exp2: ','}
		}
		c, err = it.nextToken()
		if err != nil {
			return err
		}
	}
}

// SkipObject skips an object
func (it *Iterator) SkipObject() error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	if c != '{' {
		return UnexpectedByteError{got: c, exp: '{'}
	}
	it.head++
	return skipObjectWithStack(it, c)
}
