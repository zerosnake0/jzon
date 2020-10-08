package jzon

import (
	"sync"
)

var (
	defaultIteratorPool = newIteratorPool()
)

type iteratorPool struct {
	pool sync.Pool
}

func newIteratorPool() *iteratorPool {
	return &iteratorPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &Iterator{
					tmpBuffer: make([]byte, 64),
					fixbuf:    make([]byte, 64),
				}
			},
		},
	}
}

func (p *iteratorPool) borrowIterator() *Iterator {
	it := p.pool.Get().(*Iterator)
	return it
}

func (p *iteratorPool) returnIterator(it *Iterator) {
	it.reset()
	p.pool.Put(it)
}
