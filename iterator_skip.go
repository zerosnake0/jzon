package jzon

var (
	skipFunctions [charNum]func(it *Iterator, c byte) error
)

func init() {
	skipFunctions['"'] = skipString
	skipFunctions['n'] = func(it *Iterator, _ byte) error {
		return it.expectBytes("ull")
	}
	skipFunctions['t'] = func(it *Iterator, _ byte) error {
		return it.expectBytes("rue")
	}
	skipFunctions['f'] = func(it *Iterator, _ byte) error {
		return it.expectBytes("alse")
	}
	skipFunctions['['] = skipArrayWithStack
	skipFunctions['{'] = skipObjectWithStack
	for _, c := range []byte("-0123456789") {
		skipFunctions[c] = skipNumber
	}
	errFunc := func(it *Iterator, c byte) error {
		return UnexpectedByteError{got: c}
	}
	for i := 0; i < charNum; i++ {
		if skipFunctions[i] == nil {
			skipFunctions[i] = errFunc
		}
	}
}

func skipWithStack(it *Iterator, top stackElement, s *stack) (err error) {
	var c byte
	for {
		c, err = it.nextToken()
		if err != nil {
			return err
		}
		it.head += 1
		if top&1 == 0 {
			// stackElementObjectBegin
			// stackElementObject
			if c == '}' {
				if top = s.pop(); top == stackElementNone {
					return nil
				}
				continue
			}
			if top == stackElementObject {
				if c != ',' {
					return UnexpectedByteError{got: c, exp: ','}
				}
				c, err = it.nextToken()
				if err != nil {
					return err
				}
				it.head += 1
			}
			if c != '"' {
				return UnexpectedByteError{got: c, exp: '"'}
			}
			if err = it.skipObjectField(); err != nil {
				return err
			}
			if c, err = it.nextToken(); err != nil {
				return err
			}
			it.head += 1
			switch c {
			case '[':
				s.pushObject()
				top = stackElementArrayBegin
			case '{':
				s.pushObject()
				top = stackElementObjectBegin
			default:
				if err = skipFunctions[c](it, c); err != nil {
					return err
				}
				top = stackElementObject
			}
		} else {
			// stackElementArrayBegin
			// stackElementArray
			if c == ']' {
				if top = s.pop(); top == stackElementNone {
					return nil
				}
				continue
			}
			if top == stackElementArray {
				if c != ',' {
					return UnexpectedByteError{got: c, exp: ','}
				}
				c, err = it.nextToken()
				if err != nil {
					return err
				}
				it.head += 1
			}
			switch c {
			case '[':
				s.pushArray()
				top = stackElementArrayBegin
			case '{':
				s.pushArray()
				top = stackElementObjectBegin
			default:
				if err = skipFunctions[c](it, c); err != nil {
					return err
				}
				top = stackElementArray
			}
		}
	}
}

func (it *Iterator) Skip() error {
	c, err := it.nextToken()
	if err != nil {
		return err
	}
	it.head += 1
	return skipFunctions[c](it, c)
}
