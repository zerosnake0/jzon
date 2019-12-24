package jzon

import (
	"log"
	"testing"
)

func TestValEncoder_Array(t *testing.T) {
	t.Run("pointer", func(t *testing.T) {
		arr := [...]int{1, 2, 3}
		checkEncodeValueWithStandard(t, DefaultEncoder, &arr)
	})
	t.Run("non pointer", func(t *testing.T) {
		arr := [...]int{1, 2, 3}
		checkEncodeValueWithStandard(t, DefaultEncoder, arr)
	})
	t.Run("empty", func(t *testing.T) {
		arr := [...]int{}
		checkEncodeValueWithStandard(t, DefaultEncoder, arr)
	})
	t.Run("empty pointer", func(t *testing.T) {
		arr := [...]int{}
		checkEncodeValueWithStandard(t, DefaultEncoder, &arr)
	})
	t.Run("array of pointer", func(t *testing.T) {
		i := 1
		log.Printf("%p", &i)
		arr := [...]*int{(*int)(nil), &i}
		checkEncodeValueWithStandard(t, DefaultEncoder, arr)
	})

	t.Run("one element array", func(t *testing.T) {
		arr := [...]int{1}
		checkEncodeValueWithStandard(t, DefaultEncoder, arr)
	})
	t.Run("one element array of pointer", func(t *testing.T) {
		arr := [...]*int{(*int)(nil)}
		checkEncodeValueWithStandard(t, DefaultEncoder, arr)
	})
	t.Run("one element array of pointer 2", func(t *testing.T) {
		i := 1
		arr := [...]*int{&i}
		checkEncodeValueWithStandard(t, DefaultEncoder, arr)
	})
}
