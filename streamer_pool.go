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
				return &Streamer{}
			},
		},
	}
}

func (p *StreamerPool) BorrowStreamer() *Streamer {
	s := p.pool.Get().(*Streamer)
	s.buffer = getEmptyByteSlice()
	s.poped = false
	return s
}

func (p *StreamerPool) ReturnStreamer(s *Streamer) {
	releaseByteSlice(s.buffer)
	s.buffer = nil
	p.pool.Put(s)
}
