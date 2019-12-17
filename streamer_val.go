package jzon

func (s *Streamer) Value(o interface{}) *Streamer {
	if s.Error != nil {
		return s
	}
	if o == nil {
		s.null()
		return s
	}
	// TODO:
	return s
}
