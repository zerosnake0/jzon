package jzon

type Encoder struct {
	s   *Streamer
	err error
}

func (e *Encoder) Release() {
	e.s.Release()
	e.s = nil
}

func encodeInternal(s *Streamer, v interface{}) error {
	if err := s.Value(v).Flush(); err != nil {
		return err
	}
	_, err := s.writer.Write([]byte{'\n'})
	if err != nil {
		return err
	}
	s.Reset(s.writer)
	return nil
}

func (e *Encoder) Encode(v interface{}) error {
	if e.err == nil {
		e.err = encodeInternal(e.s, v)
	}
	return e.err
}

func (e *Encoder) SetEscapeHTML(on bool) {
	e.s.EscapeHTML(on)
}

// func (e *Encoder) SetIndent(prefix, indent string) {
// 	e.s.SetIndent(prefix, indent)
// }
