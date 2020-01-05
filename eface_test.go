package jzon

import (
	"log"
	"reflect"
	"runtime/debug"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func TestEface(t *testing.T) {
	f := func(o interface{}) *eface {
		return (*eface)(unsafe.Pointer(&o))
	}
	a := 1
	b := 2
	ef1 := f(&a)
	ef2 := f(&b)
	require.Equal(t, ef1.rtype, ef2.rtype)
	// require.Equal(t, uintptr(strconv.IntSize/8), uintptr(ef2.data)-uintptr(ef1.data))
	t.Logf("%x %x", ef1.rtype, ef1.data)
	t.Logf("%x %x", ef2.rtype, ef2.data)

	// with reflect
	r := reflect.TypeOf(&a)
	ef3 := (*eface)(unsafe.Pointer(&r))
	require.Equal(t, ef1.rtype, uintptr(ef3.data))
	require.Equal(t, ef1.rtype, rtypeOfType(r))

	// pack
	packed := packEFace(rtype(ef3.data), ef1.data)
	t.Logf("%+v", packed)
	v, ok := packed.(*int)
	t.Logf("%+v %+v", v, ok)
	require.True(t, ok)
	require.Equal(t, &a, v)
}

func TestEface2(t *testing.T) {
	func() {
		i := 1
		var o interface{} = i
		ef := packEFace(rtypeOfType(reflect.TypeOf(i)), unsafe.Pointer(&i))

		t.Log("o", o)
		t.Log("ef", ef)
		i2 := ef.(int)
		t.Log("i2", i2)

		i = 2
		t.Log("o", o)
		t.Log("ef", ef)
		t.Log("i2", i2)
	}()
	debug.FreeOSMemory()
}

type testEFstruct struct {
	a int
}

func (t testEFstruct) Foo() {
	log.Printf("calling %x", unsafe.Pointer(&t))
	t.a++
}

func (t testEFstruct) Int() int {
	return t.a
}

func TestEface3(t *testing.T) {
	func() {
		var st testEFstruct
		log.Printf("&st, %p", &st)
		log.Printf("st.a, %d", st.a)
		type ifoo interface {
			Foo()
			Int() int
		}

		var foo ifoo = st
		foo.Foo()
		log.Printf("st.a, %d", st.a)
		log.Printf("foo.Int(), %d", foo.Int())
		foo.Foo()
		log.Printf("foo.Int(), %d", foo.Int())

		ef := packEFace(rtypeOfType(reflect.TypeOf(st)), unsafe.Pointer(&st))
		foo = ef.(ifoo)
		foo.Foo()
		log.Printf("st.a, %d", st.a)
	}()
	debug.FreeOSMemory()
}
