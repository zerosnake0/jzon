package jzon

import (
	"unsafe"
)

func localStringToBytes(s string) []byte {
	sh := (*stringHeader)(unsafe.Pointer(&s))
	return *(*[]byte)(unsafe.Pointer(&sliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}))
}

func localByteToString(buf []byte) string {
	return *(*string)(unsafe.Pointer(&buf))
}
