package jzon

func Unmarshal(data []byte, o interface{}) error {
	return DefaultDecoder.Unmarshal(data, o)
}

func Marshal(o interface{}) ([]byte, error) {
	return DefaultEncoder.Marshal(o)
}

func Marshal2(o interface{}) ([]byte, error) {
	return DefaultEncoder.Marshal2(o)
}

func Valid(data []byte) bool {
	it := NewIterator()
	b := it.Valid(data)
	ReturnIterator(it)
	return b
}
