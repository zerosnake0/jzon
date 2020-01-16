package jzon

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
)

// for fast reset
type iteratorEmbedded struct {
	/*
	 * The following attributes must be able to set zero by memset
	 */
	capture bool
	offset  int

	// the current index position
	head int

	// eface checkpoint
	lastEfaceOffset int
	lastEfacePtr    uintptr

	Context interface{} // custom iteration context
}

type Iterator struct {
	decoder *Decoder

	reader io.Reader
	buffer []byte
	fixbuf []byte

	// a temp buffer is needed for string reading
	// which include utf8 conversion
	tmpBuffer []byte

	// the current tail position in buffer
	tail int

	iteratorEmbedded
}

func NewIterator() *Iterator {
	return DefaultDecoder.NewIterator()
}

func ReturnIterator(it *Iterator) {
	DefaultDecoder.ReturnIterator(it)
}

func (it *Iterator) reset() {
	it.reader = nil
	it.buffer = nil
	it.tail = 0

	// fast reset
	it.iteratorEmbedded = iteratorEmbedded{}
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
	it.reader = r
	it.buffer = it.fixbuf[:cap(it.fixbuf)]
	it.tail = 0

	// fast reset
	it.iteratorEmbedded = iteratorEmbedded{}
}

func (it *Iterator) ResetBytes(data []byte) {
	it.reader = nil
	it.buffer = data
	it.tail = len(data)

	// fast reset
	it.iteratorEmbedded = iteratorEmbedded{}
}

func (it *Iterator) Buffer() []byte {
	return it.buffer[:it.tail]
}

func (it *Iterator) errorLocation() []byte {
	var (
		head int
		tail int
	)
	if it.head > 20 {
		head = it.head - 20
	}
	if it.tail-it.head < 20 {
		tail = it.tail
	} else {
		tail = it.head + 20
	}
	return it.buffer[head:tail]
}

func (it *Iterator) WrapError(err error) *DecodeError {
	if e, ok := err.(*DecodeError); ok {
		return e
	}
	return &DecodeError{
		reason:   err,
		location: string(it.errorLocation()),
	}
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
	for {
		if it.capture {
			var buf [bufferSize]byte
			n, err = it.reader.Read(buf[:])
			it.buffer = append(it.buffer[:it.tail], buf[:n]...)
			it.tail += n
			// save internal buffer for reuse
			it.fixbuf = it.buffer
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
		// n == 0 && err == nil
		// the implementation of the reader is wrong
		runtime.Gosched()
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
func (it *Iterator) nextToken() (ret byte, err error) {
	for {
		i := it.head
		for ; i < it.tail; i++ {
			c := it.buffer[i]
			if c <= ' ' {
				if valueTypeMap[c] == WhiteSpaceValue {
					continue
				}
			}
			it.head = i
			return c, nil
		}
		// the head and tail will be reset by readMore
		it.head = i
		if err := it.readMore(); err != nil {
			return 0, err
		}
	}
}

// Read until the first valid token is found, only the whitespaces are consumed
func (it *Iterator) NextValueType() (ValueType, error) {
	v, err := it.nextToken()
	return valueTypeMap[v], err
}

func (it *Iterator) Unmarshal(data []byte, obj interface{}) error {
	it.ResetBytes(data)
	err := it.ReadVal(obj)
	if err != nil {
		return err
	}
	_, err = it.nextToken()
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
	_, err = it.nextToken()
	return err == io.EOF
}
