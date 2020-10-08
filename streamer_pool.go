package jzon

import (
	"sync"
)

var (
	defaultStreamerPool = newStreamerPool()
)

type streamerPool struct {
	pool sync.Pool
}

func newStreamerPool() *streamerPool {
	return &streamerPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &Streamer{
					buffer: make([]byte, 0, 64),
				}
			},
		},
	}
}

func (p *streamerPool) borrowStreamer() *Streamer {
	s := p.pool.Get().(*Streamer)
	return s
}

func (p *streamerPool) returnStreamer(s *Streamer) {
	s.reset()
	p.pool.Put(s)
}
