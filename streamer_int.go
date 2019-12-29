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

func appendUint8(b []byte, v uint8) []byte {
	return appendLeadingDigit(b, digits[v])
}

func appendInt8(b []byte, v int8) []byte {
	if v < 0 {
		b = append(b, '-')
		return appendUint8(b, uint8(-v))
	} else {
		return appendUint8(b, uint8(v))
	}
}

func (s *Streamer) quotedInt8(v int8) *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = append(s.buffer, '"')
	s.buffer = appendInt8(s.buffer, v)
	s.buffer = append(s.buffer, '"')
	return s
}

func (s *Streamer) Int8(v int8) *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = appendInt8(s.buffer, v)
	return s
}

func (s *Streamer) quotedUint8(v uint8) *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = append(s.buffer, '"')
	s.buffer = appendUint8(s.buffer, v)
	s.buffer = append(s.buffer, '"')
	return s
}

func (s *Streamer) Uint8(v uint8) *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = appendUint8(s.buffer, v)
	return s
}

func appendUint16(b []byte, v uint16) []byte {
	q1 := v / nSmalls
	if q1 == 0 {
		return appendLeadingDigit(b, digits[v])
	}
	r1 := v - q1*nSmalls
	b = appendLeadingDigit(b, digits[q1])
	return appendDigit(b, digits[r1])
}

func appendInt16(b []byte, v int16) []byte {
	if v < 0 {
		b = append(b, '-')
		return appendUint16(b, uint16(-v))
	} else {
		return appendUint16(b, uint16(v))
	}
}

func (s *Streamer) quotedInt16(v int16) *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = append(s.buffer, '"')
	s.buffer = appendInt16(s.buffer, v)
	s.buffer = append(s.buffer, '"')
	return s
}

func (s *Streamer) Int16(v int16) *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = appendInt16(s.buffer, v)
	return s
}

func (s *Streamer) quotedUint16(v uint16) *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = append(s.buffer, '"')
	s.buffer = appendUint16(s.buffer, v)
	s.buffer = append(s.buffer, '"')
	return s
}

func (s *Streamer) Uint16(v uint16) *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = appendUint16(s.buffer, v)
	return s
}

func appendUint32(b []byte, v uint32) []byte {
	q1 := v / nSmalls
	if q1 == 0 {
		return appendLeadingDigit(b, digits[v])
	}
	r1 := v - q1*nSmalls
	q2 := q1 / nSmalls
	if q2 == 0 {
		b = appendLeadingDigit(b, digits[q1])
		return appendDigit(b, digits[r1])
	}
	r2 := q1 - q2*nSmalls
	q3 := q2 / nSmalls
	if q3 == 0 {
		b = appendLeadingDigit(b, digits[q2])
	} else {
		r3 := q2 - q3*nSmalls
		// max 10 digit for int32/uint32
		b = append(b, byte(q3+'0'))
		b = appendDigit(b, digits[r3])
	}
	b = appendDigit(b, digits[r2])
	return appendDigit(b, digits[r1])
}

func appendInt32(b []byte, v int32) []byte {
	if v < 0 {
		b = append(b, '-')
		return appendUint32(b, uint32(-v))
	} else {
		return appendUint32(b, uint32(v))
	}
}

func (s *Streamer) quotedInt32(v int32) *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = append(s.buffer, '"')
	s.buffer = appendInt32(s.buffer, v)
	s.buffer = append(s.buffer, '"')
	return s
}

func (s *Streamer) Int32(v int32) *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = appendInt32(s.buffer, v)
	return s
}

func (s *Streamer) quotedUint32(v uint32) *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = append(s.buffer, '"')
	s.buffer = appendUint32(s.buffer, v)
	s.buffer = append(s.buffer, '"')
	return s
}

func (s *Streamer) Uint32(v uint32) *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = appendUint32(s.buffer, v)
	return s
}

func appendUint64(b []byte, v uint64) []byte {
	q1 := v / nSmalls
	if q1 == 0 {
		return appendLeadingDigit(b, digits[v])
	}
	r1 := v - q1*nSmalls
	q2 := q1 / nSmalls
	if q2 == 0 {
		b = appendLeadingDigit(b, digits[q1])
	} else {
		r2 := q1 - q2*nSmalls
		q3 := q2 / nSmalls
		if q3 == 0 {
			b = appendLeadingDigit(b, digits[q2])
		} else {
			r3 := q2 - q3*nSmalls
			q4 := q3 / nSmalls
			if q4 == 0 {
				b = appendLeadingDigit(b, digits[q3])
			} else {
				r4 := q3 - q4*nSmalls
				q5 := q4 / nSmalls
				if q5 == 0 {
					b = appendLeadingDigit(b, digits[q4])
				} else {
					r5 := q4 - q5*nSmalls
					q6 := q5 / nSmalls
					if q6 == 0 {
						b = appendLeadingDigit(b, digits[q5])
					} else {
						b = appendLeadingDigit(b, digits[q6])
						r6 := q5 - q6*nSmalls
						b = appendDigit(b, digits[r6])
					}
					b = appendDigit(b, digits[r5])
				}
				b = appendDigit(b, digits[r4])
			}
			b = appendDigit(b, digits[r3])
		}
		b = appendDigit(b, digits[r2])
	}
	b = appendDigit(b, digits[r1])
	return b
}

func appendInt64(b []byte, v int64) []byte {
	if v < 0 {
		b = append(b, '-')
		return appendUint64(b, uint64(-v))
	} else {
		return appendUint64(b, uint64(v))
	}
}

func (s *Streamer) quotedInt64(v int64) *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = append(s.buffer, '"')
	s.buffer = appendInt64(s.buffer, v)
	s.buffer = append(s.buffer, '"')
	return s
}

func (s *Streamer) Int64(v int64) *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = appendInt64(s.buffer, v)
	return s
}

func (s *Streamer) quotedUint64(v uint64) *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = append(s.buffer, '"')
	s.buffer = appendUint64(s.buffer, v)
	s.buffer = append(s.buffer, '"')
	return s
}

func (s *Streamer) Uint64(v uint64) *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = appendUint64(s.buffer, v)
	return s
}

func (s *Streamer) Int(v int) *Streamer {
	return s.Int64(int64(v))
}

func (s *Streamer) Uint(v uint) *Streamer {
	return s.Uint64(uint64(v))
}
