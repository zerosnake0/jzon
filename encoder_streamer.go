package jzon

func (enc *Encoder) NewStreamer() *Streamer {
	s := defaultStreamerPool.BorrowStreamer()
	s.encoder = enc
	if enc.escapeHtml {
		s.safeSet = htmlSafeSet[:]
	} else {
		s.safeSet = safeSet[:]
	}
	return s
}

func (enc *Encoder) ReturnStreamer(s *Streamer) {
	s.encoder = nil
	defaultStreamerPool.ReturnStreamer(s)
}
