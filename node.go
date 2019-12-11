package jzon

import (
	"sync"
)

var (
	nodeStackPool = sync.Pool{
		New: func() interface{} {
			return &nodeStack{}
		},
	}
)

type readNode struct {
	m     map[string]interface{}
	s     []interface{}
	field string
}

type nodeStack struct {
	stack []readNode
}

func (ns *nodeStack) initArray(s []interface{}) *nodeStack {
	ns.stack = append(ns.stack[:0], readNode{
		s: s,
	})
	return ns
}

func (ns *nodeStack) initObject(m map[string]interface{}) *nodeStack {
	ns.stack = append(ns.stack[:0], readNode{
		m: m,
	})
	return ns
}

func (ns *nodeStack) pushArray(field string) *nodeStack {
	ns.stack = append(ns.stack, readNode{
		s:     make([]interface{}, 0),
		field: field,
	})
	return ns
}

func (ns *nodeStack) pushObject(field string) *nodeStack {
	ns.stack = append(ns.stack, readNode{
		m:     map[string]interface{}{},
		field: field,
	})
	return ns
}

func (ns *nodeStack) topObject() map[string]interface{} {
	return ns.stack[len(ns.stack)-1].m
}

func (ns *nodeStack) topArray() []interface{} {
	return ns.stack[len(ns.stack)-1].s
}

func (ns *nodeStack) setTopObject(key string, value interface{}) {
	ns.stack[len(ns.stack)-1].m[key] = value
}

func (ns *nodeStack) appendTopArray(value interface{}) {
	l := len(ns.stack) - 1
	ns.stack[l].s = append(ns.stack[l].s, value)
}

func (ns *nodeStack) popObject() {
	l := len(ns.stack) - 1
	first := &ns.stack[l]
	next := &ns.stack[l-1]
	if next.m != nil {
		next.m[first.field] = first.m
	} else {
		next.s = append(next.s, first.m)
	}
	ns.stack = ns.stack[:l]
}

func (ns *nodeStack) popArray() {
	l := len(ns.stack) - 1
	first := &ns.stack[l]
	next := &ns.stack[l-1]
	if next.m != nil {
		next.m[first.field] = first.s
	} else {
		next.s = append(next.s, first.s)
	}
	ns.stack = ns.stack[:l]
}

func releaseNodeStack(ns *nodeStack) {
	ns.stack = ns.stack[:0]
	nodeStackPool.Put(ns)
}
