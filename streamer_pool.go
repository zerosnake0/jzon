package jzon

import (
	"sync"
)

var (
	defaultStreamerPool = NewStreamerPool()
)

type StreamerPool struct {
	pool sync.Pool
}

func NewStreamerPool() *StreamerPool {
	return &StreamerPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &Streamer{
					// tmpBuffer: getByteSlice(),
				}
			},
		},
	}
}

func (p *StreamerPool) BorrowStreamer() *Streamer {
	return p.pool.Get().(*Streamer)
}

func (p *StreamerPool) ReturnStreamer(s *Streamer) {
	// it.reset()
	p.pool.Put(s)
}
