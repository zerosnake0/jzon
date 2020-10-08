package jzon

// NewStreamer returns a new streamer.
func (encCfg *EncoderConfig) NewStreamer() *Streamer {
	s := defaultStreamerPool.BorrowStreamer()
	s.cfg = encCfg
	s.EscapeHTML(s.cfg.escapeHTML)
	return s
}

func (encCfg *EncoderConfig) returnStreamer(s *Streamer) {
	s.cfg = nil
	defaultStreamerPool.ReturnStreamer(s)
}
