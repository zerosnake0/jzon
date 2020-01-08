package jzon

import (
	"reflect"
	"unsafe"
)

func localStringToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}))
}

func localByteToString(buf []byte) string {
	return *(*string)(unsafe.Pointer(&buf))
}
