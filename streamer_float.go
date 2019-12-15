package jzon

import (
	"math"
)

func (s *Streamer) Float32(f float32) *Streamer {
	if math.IsInf(float64(f), 0) {
		panic(FloatIsInfinity)
	}
	if math.IsNaN(float64(f)) {
		panic(FloatIsNan)
	}
	// TODO:
	return s
}

func (s *Streamer) Float64(f float64) *Streamer {
	if math.IsInf(f, 0) {
		panic(FloatIsInfinity)
	}
	if math.IsNaN(f) {
		panic(FloatIsNan)
	}
	// TODO:
	return s
}
