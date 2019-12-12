package jzon

import (
	"io"
)

type Streamer struct {
	encoder *Encoder

	writer io.Writer
	buffer []byte
}

func (s *Streamer) Reset(w io.Writer) {
	s.writer = w
}
