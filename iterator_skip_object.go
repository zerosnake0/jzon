package jzon

func skipObjectWithStack(it *Iterator, _ byte) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	if c == '}' {
		it.head += 1
		return nil
	}
	for {
		if c != '"' {
			return UnexpectedByteError{got: c, exp: '"'}
		}
		it.head += 1
		if err = it.skipObjectField(); err != nil {
			return err
		}
		c, err = it.nextToken()
		if err != nil {
			return err
		}
		it.head += 1
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
		it.head += 1
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

func (it *Iterator) SkipObject() error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	if c != '{' {
		return UnexpectedByteError{got: c, exp: '{'}
	}
	it.head += 1
	return skipObjectWithStack(it, c)
}
