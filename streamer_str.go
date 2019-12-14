package jzon

import (
	"unicode/utf8"
)

var (
	hex         = "0123456789abcdef"
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
	s.onVal()
	if s.encoder.escapeHtml {
		s.stringHtml(str)
	} else {
		s.string(str)
	}
	return s
}

func (s *Streamer) string(str string) {
	s.buffer = append(s.buffer, '"') // leading quote
	l := len(str)
	offset := 0
	i := 0
	for i < l {
		c := str[i]
		if c >= ' ' {
			switch c {
			case '"', '\\':
				s.buffer = append(s.buffer, str[offset:i]...)
				s.buffer = append(s.buffer, '\\', c)
				i++
				offset = i
			default:
				i++
			}
		} else {
			s.buffer = append(s.buffer, str[offset:i]...)
			switch c {
			case '\n':
				s.buffer = append(s.buffer, '\\', 'n')
			case '\r':
				s.buffer = append(s.buffer, '\\', 'r')
			case '\t':
				s.buffer = append(s.buffer, '\\', 't')
			default:
				s.buffer = append(s.buffer, '\\', 'u', '0', '0',
					hex[c>>4], hex[c&0xF])
			}
			i++
			offset = i
		}
	}
	s.buffer = append(s.buffer, str[offset:i]...)
	s.buffer = append(s.buffer, '"') // trailing quote
}

func (s *Streamer) stringHtml(str string) {
	s.buffer = append(s.buffer, '"') // leading quote
	l := len(str)
	offset := 0
	i := 0
	for i < l {
		c := str[i]
		if c < utf8.RuneSelf {
			if c >= ' ' {
				if hs := htmlSafeSet[c]; hs != "" {
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
			// TODO:
			i++
		}
	}
	s.buffer = append(s.buffer, str[offset:i]...)
	s.buffer = append(s.buffer, '"') // trailing quote
}
