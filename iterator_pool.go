package jzon

import (
	"sync"
)

var (
	defaultIteratorPool = NewIteratorPool()
)

type IteratorPool struct {
	pool sync.Pool
}

func NewIteratorPool() *IteratorPool {
	return &IteratorPool{
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

func (p *IteratorPool) BorrowIterator() *Iterator {
	it := p.pool.Get().(*Iterator)
	return it
}

func (p *IteratorPool) ReturnIterator(it *Iterator) {
	it.reset()
	p.pool.Put(it)
}
