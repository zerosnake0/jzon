package jzon

import (
	"sync"
	"sync/atomic"
)

type encoderCache = map[rtype]ValEncoder

type Encoder struct {
	cacheMu      sync.Mutex
	encoderCache atomic.Value
	tag          string
}
