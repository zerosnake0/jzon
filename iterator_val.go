package jzon

import (
	"reflect"
	"unsafe"
)

func (it *Iterator) ReadVal(obj interface{}) error {
	eface := (*eface)(unsafe.Pointer(&obj))
	if eface.data == nil {
		return ErrNilPointerReceiver
	}
	dec := it.cfg.getDecoderFromCache(eface.rtype)
	if dec == nil {
		typ := reflect.TypeOf(obj)
		if typ.Kind() != reflect.Ptr {
			return ErrPointerReceiver
		}
		dec = it.cfg.createDecoder(eface.rtype, typ)
	}
	return dec.Decode(eface.data, it, nil)
}
