package jzon

import (
	"reflect"
	"unsafe"
)

type rtype = uintptr

type eface struct {
	rtype rtype
	data  unsafe.Pointer
}

func packEFace(rtype rtype, data unsafe.Pointer) interface{} {
	var i interface{}
	e := (*eface)(unsafe.Pointer(&i))
	e.rtype = rtype
	e.data = data
	return i
}

func rtypeOfType(typ reflect.Type) rtype {
	ef := (*eface)(unsafe.Pointer(&typ))
	return rtype(ef.data)
}

type stringHeader = reflect.StringHeader
type sliceHeader = reflect.SliceHeader

/*
type stringHeader struct {
	Data unsafe.Pointer
	Len  int
}

type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}
*/
