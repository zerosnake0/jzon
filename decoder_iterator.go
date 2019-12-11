package jzon

func (dec *Decoder) NewIterator() *Iterator {
	it := defaultIteratorPool.BorrowIterator()
	it.decoder = dec
	return it
}

func (dec *Decoder) ReturnIterator(it *Iterator) {
	it.decoder = nil
	defaultIteratorPool.ReturnIterator(it)
}
