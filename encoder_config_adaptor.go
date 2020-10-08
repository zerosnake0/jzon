package jzon

func (enc *EncoderConfig) NewStreamer() *Streamer {
	s := defaultStreamerPool.BorrowStreamer()
	s.cfg = enc
	return s
}

func (enc *EncoderConfig) ReturnStreamer(s *Streamer) {
	s.cfg = nil
	defaultStreamerPool.ReturnStreamer(s)
}
