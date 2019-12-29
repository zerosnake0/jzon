package jzon

import (
	"math"
	"strconv"
)

func (s *Streamer) Float32(f float32) *Streamer {
	if s.Error != nil {
		return s
	}
	f64 := float64(f)
	if math.IsInf(f64, 0) {
		s.Error = FloatIsInfinity
		return s
	}
	if math.IsNaN(f64) {
		s.Error = FloatIsNan
		return s
	}
	s.onVal()
	fmt := byte('f')
	abs := math.Abs(f64)
	if abs != 0 {
		if f > 0 {
			if f < 1e-6 || f >= 1e21 {
				fmt = 'e'
			}
		} else {
			if f > -1e-6 || f <= -1e21 {
				fmt = 'e'
			}
		}
	}
	s.buffer = strconv.AppendFloat(s.buffer, f64, fmt, -1, 32)
	if fmt == 'e' {
		n := len(s.buffer)
		if n > 4 && s.buffer[n-4] == 'e' &&
			s.buffer[n-3] == '-' &&
			s.buffer[n-2] == '0' {
			s.buffer[n-2] = s.buffer[n-1]
			s.buffer = s.buffer[:n-1]
		}
	}
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
	s.onVal()
	fmt := byte('f')
	abs := math.Abs(f)
	if abs != 0 {
		if f > 0 {
			if f < 1e-6 || f >= 1e21 {
				fmt = 'e'
			}
		} else {
			if f > -1e-6 || f <= -1e21 {
				fmt = 'e'
			}
		}
	}
	s.buffer = strconv.AppendFloat(s.buffer, f, fmt, -1, 64)
	if fmt == 'e' {
		n := len(s.buffer)
		if n > 4 && s.buffer[n-4] == 'e' &&
			s.buffer[n-3] == '-' &&
			s.buffer[n-2] == '0' {
			s.buffer[n-2] = s.buffer[n-1]
			s.buffer = s.buffer[:n-1]
		}
	}
	return s
}
