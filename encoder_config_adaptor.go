package jzon

func (encCfg *EncoderConfig) NewStreamer() *Streamer {
	s := defaultStreamerPool.BorrowStreamer()
	s.cfg = encCfg
	return s
}

func (encCfg *EncoderConfig) ReturnStreamer(s *Streamer) {
	s.cfg = nil
	defaultStreamerPool.ReturnStreamer(s)
}
