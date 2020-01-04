package jzon

func (enc *Encoder) NewStreamer() *Streamer {
	s := defaultStreamerPool.BorrowStreamer()
	s.encoder = enc
	return s
}

func (enc *Encoder) ReturnStreamer(s *Streamer) {
	s.encoder = nil
	defaultStreamerPool.ReturnStreamer(s)
}
