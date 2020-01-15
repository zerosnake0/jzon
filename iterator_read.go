package jzon

import (
	"encoding/json"
)

var (
	readFunctions [charNum]func(it *Iterator, c byte) (interface{}, error)
)

func init() {
	readFunctions['"'] = func(it *Iterator, c byte) (interface{}, error) {
		return it.readString()
	}
	readFunctions['n'] = func(it *Iterator, c byte) (interface{}, error) {
		return nil, it.expectBytes("ull")
	}
	readFunctions['t'] = func(it *Iterator, c byte) (interface{}, error) {
		return true, it.expectBytes("rue")
	}
	readFunctions['f'] = func(it *Iterator, c byte) (interface{}, error) {
		return false, it.expectBytes("alse")
	}
	readFunctions['['] = readArrayWithStack
	readFunctions['{'] = readObjectWithStack
	for _, c := range []byte("-0123456789") {
		readFunctions[c] = func(it *Iterator, c byte) (interface{}, error) {
			if it.decoder.useNumber {
				s, err := it.readNumberAsString(c)
				return json.Number(s), err
			} else {
				return it.readFloat64(c)
			}
		}
	}
	errFunc := func(it *Iterator, c byte) (interface{}, error) {
		return nil, UnexpectedByteError{got: c}
	}
	for i := 0; i < charNum; i++ {
		if readFunctions[i] == nil {
			readFunctions[i] = errFunc
		}
	}
}

func (it *Iterator) Read() (interface{}, error) {
	c, err := it.nextToken()
	if err != nil {
		return nil, err
	}
	it.head += 1
	return readFunctions[c](it, c)
}

func readWithStack(it *Iterator, top stackElement, s *stack, ns *nodeStack) (
	_ interface{}, err error) {
	var c byte
	for {
		c, err = it.nextToken()
		if err != nil {
			return nil, err
		}
		it.head += 1
		if top&1 == 0 {
			// stackElementObjectBegin
			// stackElementObject
			if c == '}' {
				if top = s.pop(); top == stackElementNone {
					return ns.topObject(), nil
				}
				ns.popObject()
				continue
			}
			if top == stackElementObject {
				if c != ',' {
					return nil, UnexpectedByteError{got: c, exp: ','}
				}
				c, err = it.nextToken()
				if err != nil {
					return nil, err
				}
				it.head += 1
			}
			if c != '"' {
				return nil, UnexpectedByteError{got: c, exp: '"'}
			}
			field, err := it.readObjectField()
			if err != nil {
				return nil, err
			}
			c, err = it.nextToken()
			if err != nil {
				return nil, err
			}
			it.head += 1
			switch c {
			case '[':
				s.pushObject()
				ns.pushArray(field)
				top = stackElementArrayBegin
			case '{':
				s.pushObject()
				ns.pushObject(field)
				top = stackElementObjectBegin
			default:
				o, err := readFunctions[c](it, c)
				if err != nil {
					return nil, err
				}
				top = stackElementObject
				ns.setTopObject(field, o)
			}
		} else {
			// stackElementArrayBegin
			// stackElementArray
			if c == ']' {
				if top = s.pop(); top == stackElementNone {
					return ns.topArray(), nil
				}
				ns.popArray()
				continue
			}
			if top == stackElementArray {
				if c != ',' {
					return nil, UnexpectedByteError{got: c, exp: ','}
				}
				c, err = it.nextToken()
				if err != nil {
					return nil, err
				}
				it.head += 1
			}
			switch c {
			case '[':
				s.pushArray()
				ns.pushArray("")
				top = stackElementArrayBegin
			case '{':
				s.pushArray()
				ns.pushObject("")
				top = stackElementObjectBegin
			default:
				o, err := readFunctions[c](it, c)
				if err != nil {
					return nil, err
				}
				top = stackElementArray
				ns.appendTopArray(o)
			}
		}
	}
}
