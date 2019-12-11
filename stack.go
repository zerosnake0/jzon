package jzon

import (
	"sync"
)

var (
	stackPool = sync.Pool{
		New: func() interface{} {
			return &stack{
				stack: make([]uint64, 1),
			}
		},
	}
)

type stackElement = int8

const (
	stackElementNone stackElement = -1

	stackElementObjectBegin stackElement = 0 // 0b00
	stackElementObject      stackElement = 2 // 0b10

	stackElementArrayBegin stackElement = 1 // 0b01
	stackElementArray      stackElement = 3 // 0b11
)

type stack struct {
	stack []uint64
	depth uint
}

func (s *stack) init() *stack {
	s.depth = 0
	return s
}

func (s *stack) initObject() *stack {
	if len(s.stack) == 0 {
		s.stack = make([]uint64, 1)
	}
	s.stack[0] = 0
	s.depth = 1
	return s
}

func (s *stack) initArray() *stack {
	if len(s.stack) == 0 {
		s.stack = make([]uint64, 1)
	}
	s.stack[0] = 1
	s.depth = 1
	return s
}

func (s *stack) top() stackElement {
	if s.depth == 0 {
		return stackElementNone
	}
	depth := s.depth - 1
	div := depth >> 6
	mod := depth & 63
	return stackElement((s.stack[div] >> mod) & 1)
}

func (s *stack) pop() stackElement {
	if s.depth == 0 {
		return stackElementNone
	}
	depth := s.depth - 1
	div := depth >> 6
	mod := depth & 63
	s.depth -= 1
	// stackElementObjectBegin -> stackElementObject
	// stackElementArrayBegin -> stackElementArray
	return stackElement((s.stack[div]>>mod)&1) | 2
}

func (s *stack) pushObject() *stack {
	div := s.depth >> 6
	if div == uint(len(s.stack)) {
		s.stack = append(s.stack, 0)
	} else {
		s.stack[div] &= (1 << (s.depth & 63)) - 1
	}
	s.depth += 1
	return s
}

func (s *stack) pushArray() *stack {
	div := s.depth >> 6
	if div == uint(len(s.stack)) {
		s.stack = append(s.stack, 1)
	} else {
		s.stack[div] |= 1 << (s.depth & 63)
	}
	s.depth += 1
	return s
}
