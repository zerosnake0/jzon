package jzon

import (
	"reflect"
	"unsafe"
)

func localStringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&s))))
}

func localByteToString(buf []byte) string {
	return *(*string)(unsafe.Pointer(&buf))
}
