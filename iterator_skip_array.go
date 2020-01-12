package jzon

func skipArrayWithStack(it *Iterator, _ byte) error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	it.head += 1
	if c == ']' {
		return nil
	}
	for {
		switch c {
		case '{':
			s := stackPool.Get().(*stack).initArray()
			err = skipWithStack(it, stackElementObjectBegin, s)
			stackPool.Put(s)
			return err
		case '[':
			s := stackPool.Get().(*stack).initArray()
			err = skipWithStack(it, stackElementArrayBegin, s)
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
		if c == ']' {
			return nil
		}
		if c != ',' {
			return UnexpectedByteError{got: c, exp: ']', exp2: ','}
		}
		c, err = it.nextToken()
		if err != nil {
			return err
		}
		it.head += 1
	}
}

func (it *Iterator) SkipArray() error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	if c != '[' {
		return UnexpectedByteError{got: c, exp: '['}
	}
	it.head += 1
	return skipArrayWithStack(it, c)
}
