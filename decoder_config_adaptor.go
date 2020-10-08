package jzon

func (decCfg *DecoderConfig) NewIterator() *Iterator {
	it := defaultIteratorPool.BorrowIterator()
	it.cfg = decCfg
	it.useNumber = decCfg.useNumber
	it.disallowUnknownFields = decCfg.disallowUnknownFields
	return it
}

func (decCfg *DecoderConfig) ReturnIterator(it *Iterator) {
	it.cfg = nil
	defaultIteratorPool.ReturnIterator(it)
}
