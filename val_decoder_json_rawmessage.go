package jzon

import "unsafe"

type jsonRawMessageDecoder struct {
}

func (*jsonRawMessageDecoder) Decode(ptr unsafe.Pointer, it *Iterator) error {
	bytePtr := (*[]byte)(ptr)
	b, err := it.AppendRaw((*bytePtr)[:0])
	if err == nil {
		*bytePtr = b
	}
	return err
}
