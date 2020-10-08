package jzon

func (dec *DecoderConfig) NewIterator() *Iterator {
	it := defaultIteratorPool.BorrowIterator()
	it.cfg = dec
	return it
}

func (dec *DecoderConfig) ReturnIterator(it *Iterator) {
	it.cfg = nil
	defaultIteratorPool.ReturnIterator(it)
}
