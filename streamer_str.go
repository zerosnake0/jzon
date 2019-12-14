package jzon

func (s *Streamer) string(field string) {
	s.buffer = append(s.buffer, '"') // leading quote
	l := len(field)
	for i := 0; i < l; i ++ {

	}
	s.buffer = append(s.buffer, '"') // trailing quote
}
