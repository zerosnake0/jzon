package jzon

import (
	"log"
	"testing"
	"unsafe"
)

func TestValEncoder_Ptr(t *testing.T) {
	f := func(t *testing.T, o interface{}) {
		checkEncodeValueWithStandard(t, DefaultEncoder, o)
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
	t.Run("ptr", func(t *testing.T) {
		pb := (*int)(nil)
		log.Printf("&pb %p", (unsafe.Pointer)(&pb))
		log.Printf("*(&pb) %x", *(*unsafe.Pointer)((unsafe.Pointer)(&pb)))
		f(t, &pb)
	})
	t.Run("ptr 2", func(t *testing.T) {
		b := true
		log.Printf("&b %p", &b)
		pb := &b
		log.Printf("&pb %p", (unsafe.Pointer)(&pb))
		log.Printf("*(&pb) %x", *(*unsafe.Pointer)((unsafe.Pointer)(&pb)))
		f(t, &pb)
	})
}
