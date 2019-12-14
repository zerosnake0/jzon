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
				return &Iterator{}
			},
		},
	}
}

func (p *IteratorPool) BorrowIterator() *Iterator {
	it := p.pool.Get().(*Iterator)
	it.tmpBuffer = getByteSlice()
	return it
}

func (p *IteratorPool) ReturnIterator(it *Iterator) {
	it.reset()
	releaseByteSlice(it.tmpBuffer)
	it.tmpBuffer = nil
	p.pool.Put(it)
}
