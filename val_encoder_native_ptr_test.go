package jzon

import (
	"log"
	"testing"
	"unsafe"
)

func TestValEncoder_Ptr(t *testing.T) {
	f := func(t *testing.T, o interface{}) {
		checkEncodeValueWithStandard(t, DefaultEncoder, o, nil)
	}
	t.Run("nil", func(t *testing.T) {
		f(t, nil)
	})
	t.Run("true", func(t *testing.T) {
		b := true
		log.Printf("&pb %x", (unsafe.Pointer)(&b))
		log.Printf("*(&pb) %x", *(*unsafe.Pointer)((unsafe.Pointer)(&b)))
		f(t, &b)
	})
	t.Run("false", func(t *testing.T) {
		b := false
		log.Printf("&pb %p", (unsafe.Pointer)(&b))
		log.Printf("*(&pb) %x", *(*unsafe.Pointer)((unsafe.Pointer)(&b)))
		f(t, &b)
	})
	t.Run("pptr", func(t *testing.T) {
		f(t, (**bool)(nil))
	})
	t.Run("pptr 2", func(t *testing.T) {
		pb := (*bool)(nil)
		log.Printf("&pb %p", (unsafe.Pointer)(&pb))
		log.Printf("*(&pb) %x", *(*unsafe.Pointer)((unsafe.Pointer)(&pb)))
		f(t, &pb)
	})
	t.Run("pptr 3", func(t *testing.T) {
		b := true
		log.Printf("&b %p", &b)
		pb := &b
		log.Printf("&pb %p", (unsafe.Pointer)(&pb))
		log.Printf("*(&pb) %x", *(*unsafe.Pointer)((unsafe.Pointer)(&pb)))
		f(t, &pb)
	})
}
