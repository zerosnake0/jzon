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
					buffer: make([]byte, 0, 64),
				}
			},
		},
	}
}

func (p *StreamerPool) BorrowStreamer() *Streamer {
	s := p.pool.Get().(*Streamer)
	return s
}

func (p *StreamerPool) ReturnStreamer(s *Streamer) {
	s.reset()
	p.pool.Put(s)
}
