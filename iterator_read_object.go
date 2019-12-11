package jzon

func readObjectWithStack(it *Iterator, _ byte) (interface{}, error) {
	c, _, err := it.nextToken()
	if err != nil {
		return nil, err
	}
	topObj := map[string]interface{}{}
	if c == '}' {
		it.head += 1
		return topObj, nil
	}
	for {
		if c != '"' {
			return nil, UnexpectedByteError{got: c, exp: '"'}
		}
		it.head += 1
		_, field, err := it.readObjectField()
		if err != nil {
			return nil, err
		}
		c, _, err = it.nextToken()
		if err != nil {
			return nil, err
		}
		it.head += 1
		// We disabled the following switch to test benchmark
		// comparing to using builtin stack
		// the result is using ours own stack will improve the
		// performance by about 25%
		switch c {
		case '{':
			s := stackPool.Get().(*stack).initObject()
			ns := nodeStackPool.Get().(*nodeStack).
				initObject(topObj).
				pushObject(field)
			ret, err := readWithStack(it, stackElementObjectBegin, s, ns)
			releaseNodeStack(ns)
			stackPool.Put(s)
			return ret, err
		case '[':
			s := stackPool.Get().(*stack).initObject()
			ns := nodeStackPool.Get().(*nodeStack).
				initObject(topObj).
				pushArray(field)
			ret, err := readWithStack(it, stackElementArrayBegin, s, ns)
			releaseNodeStack(ns)
			stackPool.Put(s)
			return ret, err
		}
		o, err := readFunctions[c](it, c)
		if err != nil {
			return nil, err
		}
		topObj[field] = o
		c, _, err = it.nextToken()
		if err != nil {
			return nil, err
		}
		it.head += 1
		if c == '}' {
			return topObj, nil
		}
		if c != ',' {
			return nil, UnexpectedByteError{got: c, exp: '}', exp2: ','}
		}
		c, _, err = it.nextToken()
		if err != nil {
			return nil, err
		}
	}
}
