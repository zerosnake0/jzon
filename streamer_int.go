package jzon

const (
	nSmalls = 1000
)

var (
	digits [nSmalls]uint32
)

func init() {
	for i := uint32(0); i < nSmalls; i++ {
		v := ((i/100 + '0') << 16) + (((i/10)%10 + '0') << 8) + (i%10 + '0')
		if i < 10 {
			v |= 2 << 24
		} else if i < 100 {
			v |= 1 << 24
		}
		digits[i] = v
	}
}

func appendLeadingDigit(data []byte, v uint32) []byte {
	start := v >> 24
	if start == 0 {
		return append(data, byte(v>>16), byte(v>>8), byte(v))
	} else if start == 1 {
		return append(data, byte(v>>8), byte(v))
	} else { // start == 2
		return append(data, byte(v))
	}
}

func appendDigit(data []byte, v uint32) []byte {
	return append(data, byte(v>>16), byte(v>>8), byte(v))
}

func (s *Streamer) uint8(v uint8) {
	s.buffer = appendLeadingDigit(s.buffer, digits[v])
}

func (s *Streamer) Int8(v int8) *Streamer {
	s.onVal()
	if v < 0 {
		s.buffer = append(s.buffer, '-')
		s.uint8(uint8(-v))
	} else {
		s.uint8(uint8(v))
	}
	return s
}

func (s *Streamer) Uint8(v uint8) *Streamer {
	s.onVal()
	s.uint8(v)
	return s
}

func (s *Streamer) uint16(v uint16) {
	q1 := v / nSmalls
	if q1 == 0 {
		s.buffer = appendLeadingDigit(s.buffer, digits[v])
		return
	}
	r1 := v - q1*nSmalls
	s.buffer = appendLeadingDigit(s.buffer, digits[q1])
	s.buffer = appendDigit(s.buffer, digits[r1])
}

func (s *Streamer) Int16(v int16) *Streamer {
	s.onVal()
	if v < 0 {
		s.buffer = append(s.buffer, '-')
		s.uint16(uint16(-v))
	} else {
		s.uint16(uint16(v))
	}
	return s
}

func (s *Streamer) Uint16(v uint16) *Streamer {
	s.onVal()
	s.uint16(v)
	return s
}

func (s *Streamer) uint32(v uint32) {
	q1 := v / nSmalls
	if q1 == 0 {
		s.buffer = appendLeadingDigit(s.buffer, digits[v])
		return
	}
	r1 := v - q1*nSmalls
	q2 := q1 / nSmalls
	if q2 == 0 {
		s.buffer = appendLeadingDigit(s.buffer, digits[q1])
		s.buffer = appendDigit(s.buffer, digits[r1])
		return
	}
	r2 := q1 - q2*nSmalls
	q3 := q2 / nSmalls
	if q3 == 0 {
		s.buffer = appendLeadingDigit(s.buffer, digits[q2])
	} else {
		r3 := q2 - q3*nSmalls
		// max 10 digit for int32/uint32
		s.buffer = append(s.buffer, byte(q3+'0'))
		s.buffer = appendDigit(s.buffer, digits[r3])
	}
	s.buffer = appendDigit(s.buffer, digits[r2])
	s.buffer = appendDigit(s.buffer, digits[r1])
}

func (s *Streamer) Int32(v int32) *Streamer {
	s.onVal()
	if v < 0 {
		s.buffer = append(s.buffer, '-')
		s.uint32(uint32(-v))
	} else {
		s.uint32(uint32(v))
	}
	return s
}

func (s *Streamer) Uint32(v uint32) *Streamer {
	s.onVal()
	s.uint32(v)
	return s
}

func (s *Streamer) uint64(v uint64) {
	q1 := v / nSmalls
	if q1 == 0 {
		s.buffer = appendLeadingDigit(s.buffer, digits[v])
		return
	}
	r1 := v - q1*nSmalls
	q2 := q1 / nSmalls
	if q2 == 0 {
		s.buffer = appendLeadingDigit(s.buffer, digits[q1])
	} else {
		r2 := q1 - q2*nSmalls
		q3 := q2 / nSmalls
		if q3 == 0 {
			s.buffer = appendLeadingDigit(s.buffer, digits[q2])
		} else {
			r3 := q2 - q3*nSmalls
			q4 := q3 / nSmalls
			if q4 == 0 {
				s.buffer = appendLeadingDigit(s.buffer, digits[q3])
			} else {
				r4 := q3 - q4*nSmalls
				q5 := q4 / nSmalls
				if q5 == 0 {
					s.buffer = appendLeadingDigit(s.buffer, digits[q4])
				} else {
					r5 := q4 - q5*nSmalls
					q6 := q5 / nSmalls
					if q6 == 0 {
						s.buffer = appendLeadingDigit(s.buffer, digits[q5])
					} else {
						s.buffer = appendLeadingDigit(s.buffer, digits[q6])
						r6 := q5 - q6*nSmalls
						s.buffer = appendDigit(s.buffer, digits[r6])
					}
					s.buffer = appendDigit(s.buffer, digits[r5])
				}
				s.buffer = appendDigit(s.buffer, digits[r4])
			}
			s.buffer = appendDigit(s.buffer, digits[r3])
		}
		s.buffer = appendDigit(s.buffer, digits[r2])
	}
	s.buffer = appendDigit(s.buffer, digits[r1])
}

func (s *Streamer) Int64(v int64) *Streamer {
	s.onVal()
	if v < 0 {
		s.buffer = append(s.buffer, '-')
		s.uint64(uint64(-v))
	} else {
		s.uint64(uint64(v))
	}
	return s
}

func (s *Streamer) Uint64(v uint64) *Streamer {
	s.onVal()
	s.uint64(v)
	return s
}
