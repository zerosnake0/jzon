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
	s.buffer = encodeString(s.buffer, str, s.encoder.safeSet)
	return s
}

func encodeString(buffer []byte, str string, safeSet []string) []byte {
	buffer = append(buffer, '"') // leading quote
	l := len(str)
	offset := 0
	i := 0
	for i < l {
		c := str[i]
		if c < utf8.RuneSelf {
			if c >= ' ' {
				if hs := safeSet[c]; hs != "" {
					buffer = append(buffer, str[offset:i]...)
					buffer = append(buffer, hs...)
					i++
					offset = i
				} else {
					i++
				}
			} else {
				buffer = append(buffer, str[offset:i]...)
				if hs := htmlSafeSet[c]; hs != "" {
					buffer = append(buffer, hs...)
				} else {
					buffer = append(buffer, '\\', 'u', '0', '0',
						hex[c>>4], hex[c&0xF])
				}
				i++
				offset = i
			}
		} else { // c >= 0x80
			r, size := utf8.DecodeRuneInString(str[i:])
			if r == utf8.RuneError && size == 1 {
				buffer = append(buffer, str[offset:i]...)
				buffer = append(buffer, '\\', 'u',
					'f', 'f', 'f', 'd')
				i += 1
				offset = i
				continue
			}

			if r == '\u2028' || r == '\u2029' {
				buffer = append(buffer, str[offset:i]...)
				buffer = append(buffer, '\\', 'u',
					'2', '0', '2', hex[r&0xF])
				i += size
				offset = i
			} else {
				i += size
			}
		}
	}
	buffer = append(buffer, str[offset:i]...)
	buffer = append(buffer, '"') // trailing quote
	return buffer
}
