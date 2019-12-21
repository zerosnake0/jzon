package jzon

import (
	"unicode/utf8"
)

var (
	hex     = "0123456789abcdef"
	safeSet = [utf8.RuneSelf]string{
		// < ' '
		'\n': `\n`,
		'\r': `\r`,
		'\t': `\t`,
		// >= ' '
		'"':  `\"`,
		'\\': `\\`,
	}
	htmlSafeSet = [utf8.RuneSelf]string{
		// < ' '
		'\n': `\n`,
		'\r': `\r`,
		'\t': `\t`,
		// >= ' '
		'"':  `\"`,
		'\\': `\\`,
		'<':  `\u003c`,
		'>':  `\u003e`,
		'&':  `\u0026`,
	}
)

func (s *Streamer) String(str string) *Streamer {
	if s.Error != nil {
		return s
	}
	s.onVal()
	s.string(str)
	return s
}

func (s *Streamer) string(str string) {
	s.buffer = append(s.buffer, '"') // leading quote
	l := len(str)
	offset := 0
	i := 0
	for i < l {
		c := str[i]
		if c < utf8.RuneSelf {
			if c >= ' ' {
				if hs := s.safeSet[c]; hs != "" {
					s.buffer = append(s.buffer, str[offset:i]...)
					s.buffer = append(s.buffer, hs...)
					i++
					offset = i
				} else {
					i++
				}
			} else {
				s.buffer = append(s.buffer, str[offset:i]...)
				if hs := htmlSafeSet[c]; hs != "" {
					s.buffer = append(s.buffer, hs...)
				} else {
					s.buffer = append(s.buffer, '\\', 'u', '0', '0',
						hex[c>>4], hex[c&0xF])
				}
				i++
				offset = i
			}
		} else { // c >= 0x80
			s.buffer = append(s.buffer, str[offset:i]...)
			r, size := utf8.DecodeRuneInString(str[i:])
			if r == utf8.RuneError {
				// we must have size == 1 here
				// because the input is not empty
				s.buffer = append(s.buffer, '\\', 'u',
					'f', 'f', 'f', 'd')
				i += 1
				offset = i
			} else if r == '\u2028' || r == '\u2029' {
				s.buffer = append(s.buffer, '\\', 'u',
					'2', '0', '2', hex[r&0xF])
				i += size
				offset = i
			} else {
				i += size
			}
		}
	}
	s.buffer = append(s.buffer, str[offset:i]...)
	s.buffer = append(s.buffer, '"') // trailing quote
}
