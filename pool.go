package jzon

import (
	"sync"
)

const bufferSize = 64

var (
	bytesPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, bufferSize)
		},
	}
)

func getByteSlice() []byte {
	return bytesPool.Get().([]byte)
}

func getEmptyByteSlice() []byte {
	b := bytesPool.Get().([]byte)
	return b[:0]
}

func getFullByteSlice() []byte {
	b := bytesPool.Get().([]byte)
	return b[:cap(b)]
}

func releaseByteSlice(b []byte) {
	if cap(b) != 0 {
		bytesPool.Put(b)
	}
}
