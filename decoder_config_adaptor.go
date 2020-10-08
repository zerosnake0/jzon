package jzon

// NewIterator returns a new iterator.
func (decCfg *DecoderConfig) NewIterator() *Iterator {
	it := defaultIteratorPool.borrowIterator()
	it.cfg = decCfg
	it.useNumber = decCfg.useNumber
	it.disallowUnknownFields = decCfg.disallowUnknownFields
	return it
}

func (decCfg *DecoderConfig) returnIterator(it *Iterator) {
	it.cfg = nil
	defaultIteratorPool.returnIterator(it)
}
