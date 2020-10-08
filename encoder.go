package jzon

// Encoder is almost standard library compatible
// The following standard methods are not implemented
// - SetIndent
type Encoder struct {
	s   *Streamer
	err error
}

// Release encoder, encoder should not be reused after call
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

// Encode writes the JSON encoding of v to the stream,
// followed by a newline character.
func (e *Encoder) Encode(v interface{}) error {
	if e.err == nil {
		e.err = encodeInternal(e.s, v)
	}
	return e.err
}

// SetEscapeHTML specifies whether problematic HTML characters
// should be escaped inside JSON quoted strings.
// The default behavior is to escape &, <, and > to \u0026, \u003c, and \u003e
// to avoid certain safety problems that can arise when embedding JSON in HTML.
//
// In non-HTML settings where the escaping interferes with the readability
// of the output, SetEscapeHTML(false) disables this behavior.
func (e *Encoder) SetEscapeHTML(on bool) {
	e.s.EscapeHTML(on)
}

// func (e *Encoder) SetIndent(prefix, indent string) {
// 	e.s.SetIndent(prefix, indent)
// }
