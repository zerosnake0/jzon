package jzon

import (
	"bytes"
	"fmt"
	"io"
)

type Iterator struct {
	decoder *Decoder

	reader io.Reader
	buffer []byte

	// a temp buffer is needed for string reading
	// which include utf8 conversion
	tmpBuffer []byte

	capture bool
	offset  int

	// the current index position
	head int
	tail int

	// path string

	// Error error
}

func NewIterator() *Iterator {
	return DefaultDecoder.NewIterator()
}

func ReturnIterator(it *Iterator) {
	DefaultDecoder.ReturnIterator(it)
}

func (it *Iterator) reset() {
	if it.reader == nil {
		it.buffer = nil
	} else { // it.reader != nil
		it.reader = nil
		releaseByteSlice(it.buffer)
		it.buffer = nil
	}
}

/*
 * In reset methods, explicit assignment is faster than then following
 *   *it = Iterator{ ... }
 * When the above code is used, runtime.duffcopy and runtime.duffzero will be used
 * which will slow down our code (correct me if I am wrong)
 */
func (it *Iterator) Reset(r io.Reader) {
	switch v := r.(type) {
	case nil:
		it.reset()
		return
	case *bytes.Buffer:
		it.ResetBytes(v.Bytes())
		return
	}
	var b []byte
	if it.reader == nil {
		b = getFullByteSlice()
	} else {
		b = it.buffer
	}
	it.reader = r
	it.buffer = b
	it.offset = 0
	it.head = 0
	it.tail = 0
}

func (it *Iterator) ResetBytes(data []byte) {
	if it.reader != nil && it.buffer != nil {
		releaseByteSlice(it.buffer)
	}
	it.reader = nil
	it.buffer = data
	it.offset = 0
	it.head = 0
	it.tail = len(data)
}

func (it *Iterator) Buffer() []byte {
	return it.buffer[:it.tail]
}

// make sure that it.head == it.tail before call
func (it *Iterator) readMore() error {
	if it.reader == nil {
		return io.EOF
	}
	var (
		n   int
		err error
	)
	// TODO: risk of infinite loop?
	for {
		if it.capture {
			var buf [bufferSize]byte
			n, err = it.reader.Read(buf[:])
			it.buffer = append(it.buffer[:it.tail], buf[:n]...)
			it.tail += n
		} else {
			if it.head != it.tail { // debug, to be removed
				panic(fmt.Errorf("head %d, tail %d", it.head, it.tail))
			}
			n, err = it.reader.Read(it.buffer)
			it.offset += it.tail
			it.head = 0
			it.tail = n
		}
		if err != nil {
			if err == io.EOF && n > 0 {
				return nil
			}
			return err
		}
		if n > 0 {
			return nil
		}
	}
}

// will NOT skip whitespaces
// will NOT consume the character
// will report error on EOF
func (it *Iterator) nextByte() (ret byte, err error) {
	if it.head == it.tail {
		if err = it.readMore(); err != nil {
			return
		}
	}
	return it.buffer[it.head], nil
}

// will consume the characters
func (it *Iterator) expectBytes(s string) error {
	last := len(s) - 1
	j := 0
	for {
		i := it.head
		for ; i < it.tail; i++ {
			c := it.buffer[i]
			if c != s[j] {
				return UnexpectedByteError{exp: s[j], got: c}
			}
			if j == last {
				it.head = i + 1
				return nil
			}
			j++
		}
		it.head = i
		if err := it.readMore(); err != nil {
			return err
		}
	}
}

// Read until the first valid token is found, only the whitespaces are consumed
func (it *Iterator) nextToken() (ret byte, vt ValueType, err error) {
	for {
		i := it.head
		for ; i < it.tail; i++ {
			c := it.buffer[i]
			vt := valueTypeMap[c]
			if vt == WhiteSpaceValue {
				continue
			}
			it.head = i
			return c, vt, nil
		}
		// the head and tail will be reset by readMore
		it.head = i
		if err := it.readMore(); err != nil {
			return 0, InvalidValue, err
		}
	}
}

// Read until the first valid token is found, only the whitespaces are consumed
func (it *Iterator) NextValueType() (vt ValueType, err error) {
	_, vt, err = it.nextToken()
	return
}

func (it *Iterator) Unmarshal(data []byte, obj interface{}) error {
	it.ResetBytes(data)
	err := it.ReadVal(obj)
	if err != nil {
		return err
	}
	_, _, err = it.nextToken()
	if err == nil {
		return DataRemainedError
	}
	if err != io.EOF {
		return err
	}
	return nil
}

func (it *Iterator) Valid(data []byte) bool {
	it.ResetBytes(data)
	err := it.Skip()
	if err != nil {
		return false
	}
	_, _, err = it.nextToken()
	return err == io.EOF
}
