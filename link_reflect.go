package jzon

import (
	"reflect"
	"unsafe"
)

/*
 * WARNING:
 * The linked functions in this file should be used with EXTREMELY careful
 */

//go:linkname unsafe_New reflect.unsafe_New
func unsafe_New(rtype rtype) unsafe.Pointer

//go:linkname typedmemclrpartial reflect.typedmemclrpartial
func typedmemclrpartial(t rtype, ptr unsafe.Pointer, off, size uintptr)

//go:linkname unsafe_NewArray reflect.unsafe_NewArray
func unsafe_NewArray(rtype rtype, length int) unsafe.Pointer

//go:linkname typedslicecopy reflect.typedslicecopy
func typedslicecopy(rtyp rtype, dst, src reflect.SliceHeader) int

//go:linkname makemap reflect.makemap
func makemap(rtype rtype, cap int) unsafe.Pointer

//go:linkname typedmemmove reflect.typedmemmove
func typedmemmove(rtype rtype, dst, src unsafe.Pointer)

//go:linkname mapassign reflect.mapassign
//go:noescape
func mapassign(t rtype, m, key, val unsafe.Pointer)

//go:linkname maplen reflect.maplen
//go:noescape
func maplen(m unsafe.Pointer) int

//go:linkname ifaceIndir reflect.ifaceIndir
func ifaceIndir(t rtype) bool

type hiter struct {
	key   unsafe.Pointer
	value unsafe.Pointer
}

//go:noescape
//go:linkname mapiterinit reflect.mapiterinit
func mapiterinit(rtype rtype, m unsafe.Pointer) *hiter

//go:noescape
//go:linkname mapiternext reflect.mapiternext
func mapiternext(it *hiter)

func unsafeMakeSlice(rtype rtype, length, cap int) unsafe.Pointer {
	return unsafe.Pointer(&reflect.SliceHeader{
		// TODO: is this safe?
		Data: uintptr(unsafe_NewArray(rtype, cap)),
		Len:  length,
		Cap:  cap,
	})
}

func unsafeMakeMap(rtype rtype, cap int) unsafe.Pointer {
	m := makemap(rtype, cap)
	return unsafe.Pointer(&m)
}

// see reflect.add
func add(p unsafe.Pointer, x uintptr, whySafe string) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

// see reflect.grow
func unsafeGrowSlice(rtype, elemRType rtype, ptr unsafe.Pointer, newLength int) unsafe.Pointer {
	sh := (*reflect.SliceHeader)(ptr)
	if newLength < sh.Cap {
		sh.Len = newLength
		return ptr
	}
	newCap := sh.Cap
	if sh.Cap == 0 {
		newCap = newLength
	} else {
		for newCap < newLength {
			if newCap < 1024 {
				newCap <<= 1
			} else {
				newCap += newCap >> 2
			}
		}
	}
	newHeader := (*reflect.SliceHeader)(unsafeMakeSlice(rtype, newLength, newCap))
	typedslicecopy(elemRType, *newHeader, *sh)
	return unsafe.Pointer(newHeader)
}

func unsafeSliceChildPtr(ptr unsafe.Pointer, elemSize uintptr, index int) unsafe.Pointer {
	sh := (*reflect.SliceHeader)(ptr)
	return add(unsafe.Pointer(sh.Data), uintptr(index)*elemSize, "index < len")
}
