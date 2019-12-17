package jzon

import (
	"io"
)

type Streamer struct {
	encoder *Encoder

	writer io.Writer
	buffer []byte

	Error   error
	poped   bool
	safeSet []string
}

func NewStreamer() *Streamer {
	return DefaultEncoder.NewStreamer()
}

func ReturnStreamer(s *Streamer) {
	DefaultEncoder.ReturnStreamer(s)
}

func (s *Streamer) reset() {
	s.writer = nil
	s.Error = nil
}

func (s *Streamer) Reset(w io.Writer) {
	s.writer = w
}

func (s *Streamer) Flush() error {
	if s.Error != nil {
		return s.Error
	}
	if s.writer == nil {
		return NoWriterAttachedError
	}
	l := len(s.buffer)
	// see comment of io.Writer
	n, err := s.writer.Write(s.buffer)
	if n < l {
		// TODO: shall we accept the writers which
		// TODO: do not implement Write method correctly?
		copy(s.buffer, s.buffer[n:])
		s.buffer = s.buffer[:l-n]
	} else {
		s.buffer = s.buffer[:0]
	}
	return err
}

func (s *Streamer) onVal() {
	if s.poped {
		s.buffer = append(s.buffer, ',')
	} else {
		s.poped = true
	}
}

func (s *Streamer) RawString(raw string) *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = append(s.buffer, raw...)
	return s
}

func (s *Streamer) Raw(raw []byte) *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = append(s.buffer, raw...)
	return s
}

func (s *Streamer) null() {
	s.onVal()
	s.buffer = append(s.buffer, 'n', 'u', 'l', 'l')
}

func (s *Streamer) Null() *Streamer {
	if s.Error != nil {
		return s
	}
	s.null()
	return s
}

func (s *Streamer) True() *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = append(s.buffer, 't', 'r', 'u', 'e')
	return s
}

func (s *Streamer) False() *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = append(s.buffer, 'f', 'a', 'l', 's', 'e')
	return s
}

func (s *Streamer) Bool(b bool) *Streamer {
	if b {
		return s.True()
	} else {
		return s.False()
	}
}

func (s *Streamer) ObjectStart() *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = append(s.buffer, '{')
	s.poped = false
	return s
}

func (s *Streamer) Field(field string) *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.string(field)
	s.buffer = append(s.buffer, ':')
	s.poped = false
	return s
}

func (s *Streamer) ObjectEnd() *Streamer {
	if s.Error != nil {
		return s
	}
	s.buffer = append(s.buffer, '}')
	s.poped = true
	return s
}

func (s *Streamer) ArrayStart() *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.buffer = append(s.buffer, '[')
	s.poped = false
	return s
}

func (s *Streamer) ArrayEnd() *Streamer {
	if s.Error != nil {
		return s
	}
	s.buffer = append(s.buffer, ']')
	s.poped = true
	return s
}
