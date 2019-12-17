package jzon

import (
	"math"
)

func (s *Streamer) Float32(f float32) *Streamer {
	if s.Error != nil {
		return s
	}
	if math.IsInf(float64(f), 0) {
		s.Error = FloatIsInfinity
		return s
	}
	if math.IsNaN(float64(f)) {
		s.Error = FloatIsNan
		return s
	}
	// TODO:
	return s
}

func (s *Streamer) Float64(f float64) *Streamer {
	if s.Error != nil {
		return s
	}
	if math.IsInf(f, 0) {
		s.Error = FloatIsInfinity
		return s
	}
	if math.IsNaN(f) {
		s.Error = FloatIsNan
		return s
	}
	// TODO:
	return s
}
