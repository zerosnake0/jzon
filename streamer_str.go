package jzon

var (
	hex = "0123456789abcdef"
)

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
