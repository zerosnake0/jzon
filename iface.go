package jzon

import (
	"unsafe"
)

type iface struct {
	itab *itab
	data unsafe.Pointer
}

type itab struct {
	ignore uintptr
	rtype  rtype
}

func packIFace(ptr unsafe.Pointer) interface{} {
	iface := (*iface)(ptr)
	if iface.itab == nil {
		return nil
	}
	return packEFace(iface.itab.rtype, iface.data)
}
