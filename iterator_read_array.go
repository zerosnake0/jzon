package jzon

func readArrayWithStack(it *Iterator, _ byte) (interface{}, error) {
	c, err := it.nextToken()
	if err != nil {
		return nil, err
	}
	it.head += 1
	topObj := make([]interface{}, 0)
	if c == ']' {
		return topObj, nil
	}
	for {
		// We disabled the following switch to test benchmark
		// comparing to using builtin stack
		// the result is using our own stack will improve the
		// performance by about 25%
		switch c {
		case '{':
			s := stackPool.Get().(*stack).initArray()
			ns := nodeStackPool.Get().(*nodeStack).
				initArray(topObj).
				pushObject("")
			ret, err := readWithStack(it, stackElementObjectBegin, s, ns)
			releaseNodeStack(ns)
			stackPool.Put(s)
			return ret, err
		case '[':
			s := stackPool.Get().(*stack).initArray()
			ns := nodeStackPool.Get().(*nodeStack).
				initArray(topObj).
				pushArray("")
			ret, err := readWithStack(it, stackElementArrayBegin, s, ns)
			releaseNodeStack(ns)
			stackPool.Put(s)
			return ret, err
		}
		o, err := readFunctions[c](it, c)
		if err != nil {
			return nil, err
		}
		topObj = append(topObj, o)
		c, err = it.nextToken()
		if err != nil {
			return nil, err
		}
		it.head += 1
		if c == ']' {
			return topObj, nil
		}
		if c != ',' {
			return nil, UnexpectedByteError{got: c, exp: ',', exp2: ']'}
		}
		c, err = it.nextToken()
		if err != nil {
			return nil, err
		}
		it.head += 1
	}
}
