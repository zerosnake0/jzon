//+build !go1.18

package jzon

import (
	"unsafe"
)

//go:noescape
//go:linkname mapiterinit reflect.mapiterinit
func mapiterinit(rtype rtype, m unsafe.Pointer) *hiter