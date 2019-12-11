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
					tmpBuffer: getByteSlice(),
				}
			},
		},
	}
}

func (p *IteratorPool) BorrowIterator() *Iterator {
	return p.pool.Get().(*Iterator)
}

func (p *IteratorPool) ReturnIterator(it *Iterator) {
	it.reset()
	p.pool.Put(it)
}
