package jzon

import (
	"reflect"
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
