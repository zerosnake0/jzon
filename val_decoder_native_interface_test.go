package jzon

import (
	"errors"
	"io"
	"reflect"
	"runtime/debug"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

type testIface interface {
	Foo()
}

type testIfaceImpl struct {
	field string
}

func (testIfaceImpl) Foo() {
}

func copyValue(t *testing.T, in interface{}) (out interface{}) {
	if in == nil {
		return nil
	}
	typ := reflect.TypeOf(in)
	switch k := typ.Kind(); k {
	case reflect.String,
		reflect.Int,
		reflect.Uint8:
		return in
	case reflect.Map:
		v := reflect.ValueOf(in)
		newV := reflect.MakeMap(typ)
		if !v.IsNil() {
			iter := v.MapRange()
			for iter.Next() {
				k := copyValue(t, iter.Key().Interface())
				v := copyValue(t, iter.Value().Interface())
				newV.SetMapIndex(reflect.ValueOf(k),
					reflect.ValueOf(v))
			}
		}
		return newV.Interface()
	case reflect.Ptr:
		ptrValue := reflect.ValueOf(in)
		if ptrValue.IsNil() {
			newV := reflect.NewAt(typ.Elem(), nil)
			return newV.Interface()
		}
		elem := ptrValue.Elem()
		copied := copyValue(t, elem.Interface())
		newV := reflect.New(elem.Type())
		if copied != nil {
			newV.Elem().Set(reflect.ValueOf(copied))
		}
		return newV.Interface()
	case reflect.Struct:
		oldV := reflect.ValueOf(in)
		newV := reflect.New(typ).Elem()
		for i := 0; i < oldV.NumField(); i++ {
			field := oldV.Field(i)
			if field.CanInterface() {
				copiedV := copyValue(t, field.Interface())
				newV.Field(i).Set(reflect.ValueOf(copiedV))
			}
		}
		return newV.Interface()
	case reflect.Slice:
		v := reflect.ValueOf(in)
		l := v.Len()
		newV := reflect.MakeSlice(typ, l, l)
		for i := 0; i < l; i++ {
			copied := copyValue(t, v.Index(i).Interface())
			newV.Index(i).Set(reflect.ValueOf(copied))
		}
		return newV.Interface()
	case reflect.Interface:
		v := reflect.ValueOf(in)
		newV := reflect.New(typ)
		if !v.IsNil() {
			copiedV := v.Elem().Interface()
			newV.Elem().Set(reflect.ValueOf(copiedV))
		}
		return newV.Elem()
	default:
		t.Fatalf("%v(%s) not supported", typ, k.String())
		panic("should not reach here")
	}
}

func TestValDecoder_Native_Interface(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		t.Log(">>>>> initValues >>>>>")
		printValue(t, "p1", p1)
		printValue(t, "p2", p2)
		t.Log(">>>>>>>>>>>>>>>>>>>>>>")
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
		t.Log("<<<<< initValues <<<<<")
		printValue(t, "p1", p1)
		printValue(t, "p2", p2)
		t.Log("<<<<<<<<<<<<<<<<<<<<<<")
	}
	f2 := func(t *testing.T, data string, ex error, initValues ...interface{}) {
		var v1 interface{}
		var v2 interface{}
		var p1 *interface{}
		var p2 *interface{}
		if len(initValues) != 0 {
			if len(initValues) == 1 {
				v1 = copyValue(t, initValues[0])
				v2 = copyValue(t, initValues[0])
			} else {
				v1 = initValues[0]
				v2 = initValues[1]
			}
			require.Equal(t, v1, v2)
			p1 = &v1
			p2 = &v2
		}
		f(t, data, ex, p1, p2)
	}
	f3 := func(t *testing.T, data string, ex error) {
		f2(t, data, ex, "dummy")
	}

	// eface
	t.Run("nil pointer", func(t *testing.T) {
		f2(t, "null", NilPointerReceiverError)
	})
	t.Run("eof", func(t *testing.T) {
		f3(t, "", io.EOF)
	})
	t.Run("invalid first byte", func(t *testing.T) {
		f2(t, `+`, UnexpectedByteError{}, nil)
	})
	t.Run("nil init value", func(t *testing.T) {
		f2(t, `{"a":"b"}`, nil, nil)
	})
	t.Run("nil typed pointer 1", func(t *testing.T) {
		f2(t, `{"a":"b"}`, nil, (*int)(nil))
	})
	t.Run("nil typed pointer 2", func(t *testing.T) {
		f2(t, `{"a":"b"}`, nil, (**int)(nil))
	})
	t.Run("nil typed pointer 3", func(t *testing.T) {
		f2(t, `{"a":"b"}`, nil, (***int)(nil))
	})
	t.Run("non compatible value", func(t *testing.T) {
		f2(t, `{"a":"b"}`, nil, 1)
	})
	t.Run("non compatible value 2", func(t *testing.T) {
		f2(t, `{"a":"b"}`, nil, testJsonUnmarshaler{
			data: "123",
			err:  errors.New("test"),
		}, testJsonUnmarshaler{
			data: "123",
			err:  errors.New("test"),
		})
	})
	t.Run("eof 2", func(t *testing.T) {
		f2(t, ``, io.EOF, &testJsonUnmarshaler{
			data: "123",
			err:  errors.New("test error"),
		}, &testJsonUnmarshaler{
			data: "123",
			err:  errors.New("test error"),
		})
	})
	t.Run("non compatible value 3", func(t *testing.T) {
		ex := errors.New("test error")
		f2(t, `{"a":"b"}`, ex,
			&testJsonUnmarshaler{
				data: "123",
				err:  ex,
			}, &testJsonUnmarshaler{
				data: "123",
				err:  ex,
			})
	})
	t.Run("null on non nil init value", func(t *testing.T) {
		f2(t, `null`, nil, "1")
	})
	t.Run("different type with non nil init value", func(t *testing.T) {
		f2(t, `123`, nil, "1")
	})
	t.Run("null with nil pointer init value", func(t *testing.T) {
		f2(t, `null`, nil, (*int)(nil))
	})
	t.Run("null with nil pointer init value 2", func(t *testing.T) {
		f2(t, `null`, nil, (**int)(nil))
	})
	t.Run("invalid null with non nil pointer init value", func(t *testing.T) {
		i := 1
		f2(t, `nul`, io.EOF, &i)
	})
	t.Run("null with non nil pointer init value", func(t *testing.T) {
		i := 1
		f2(t, `null`, nil, &i)
	})
	t.Run("non null with non nil pointer init value", func(t *testing.T) {
		i := 1
		f2(t, `"test"`, InvalidDigitError{}, &i)
	})
	t.Run("null with pt", func(t *testing.T) {
		var v interface{}
		pv := &v                  // *interface{}
		var ppv interface{} = &pv // **interface{}
		f2(t, `null`, nil, ppv)
	})
	t.Run("null with ptr 2", func(t *testing.T) {
		// var v interface{}
		// pv := &v // *interface{}
		// var ppv interface{} = &pv // **interface{}
		f2(t, `null`, nil, (**interface{})(nil))
	})
	t.Run("non null with ptr 2", func(t *testing.T) {
		f2(t, `"test"`, nil, (**interface{})(nil))
	})
	t.Run("null with ptr 3-1", func(t *testing.T) {
		var v interface{}
		pv := &v                // *interface{}
		f2(t, `null`, nil, &pv) // **interface{}
	})
	t.Run("non null with ptr 3-1", func(t *testing.T) {
		var v interface{}
		pv := &v              // *interface{}
		f2(t, `24`, nil, &pv) // **interface{}
	})
	t.Run("null with ptr 3-2", func(t *testing.T) {
		var v interface{}
		pv := &v
		ppv := &pv
		var pppv interface{} = &ppv
		f2(t, `null`, nil, &pppv)
	})
	t.Run("non null with ptr 3-2", func(t *testing.T) {
		var v interface{}
		pv := &v
		ppv := &pv
		var pppv interface{} = &ppv
		f2(t, `"test"`, nil, &pppv)
	})

	// iface
	t.Run("iface eof", func(t *testing.T) {
		var um1 testIface
		var um2 testIface
		f(t, ``, io.EOF, &um1, &um2)
	})
	t.Run("iface invalid null", func(t *testing.T) {
		var um1 testIface
		var um2 testIface
		f(t, `nul`, io.EOF, &um1, &um2)
	})
	t.Run("iface null 1", func(t *testing.T) {
		var um1 testIface
		var um2 testIface
		f(t, `null`, nil, &um1, &um2)
	})
	t.Run("iface null 2", func(t *testing.T) {
		var um1 testIface = testIfaceImpl{
			field: "test",
		}
		var um2 testIface = testIfaceImpl{
			field: "test",
		}
		f(t, `null`, nil, &um1, &um2)
	})
	t.Run("iface null 3", func(t *testing.T) {
		var um1 testIface = &testIfaceImpl{
			field: "test",
		}
		var um2 testIface = &testIfaceImpl{
			field: "test",
		}
		f(t, `null`, nil, &um1, &um2)
	})
	t.Run("iface not null 1", func(t *testing.T) {
		var um1 testIface
		var um2 testIface
		f(t, `{}`, IFaceError, &um1, &um2)
	})
	t.Run("iface not null 2", func(t *testing.T) {
		var um1 testIface = testIfaceImpl{
			field: "test",
		}
		var um2 testIface = testIfaceImpl{
			field: "test",
		}
		f(t, `{}`, PointerReceiverError, &um1, &um2)
	})
	t.Run("iface not null 3", func(t *testing.T) {
		var um1 testIface = &testIfaceImpl{
			field: "test",
		}
		var um2 testIface = &testIfaceImpl{
			field: "test",
		}
		f(t, `{}`, nil, &um1, &um2)
	})
	debug.FreeOSMemory()
}

type testIfaceDecoder struct {
}

func (*testIfaceDecoder) Decode(ptr unsafe.Pointer, it *Iterator, _ *DecOpts) error {
	o, err := it.Read()
	if err != nil {
		return err
	}
	*(*interface{})(ptr) = o
	return nil
}

func TestValDecoder_Native_Interface_Loop(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		t.Log(">>>>> initValues >>>>>")
		printValue(t, "p1", p1)
		printValue(t, "p2", p2)
		t.Log(">>>>>>>>>>>>>>>>>>>>>>")
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
		t.Log("<<<<< initValues <<<<<")
		printValue(t, "p1", p1)
		printValue(t, "p2", p2)
		t.Log("<<<<<<<<<<<<<<<<<<<<<<")
	}
	/*
	 * See comment of efaceDecoder
	 */
	t.Run("compatible with standard", func(t *testing.T) {
		// we have different behavior with standard library for this test
		skipTest(t, "eface looping")
		var o1 interface{}
		o1 = &o1
		var o2 interface{}
		o2 = &o2
		f(t, `1`, nil, o1, o2)
	})
	t.Run("loop 1", func(t *testing.T) {
		type iface interface{}
		var o1 iface
		o1 = &o1
		err := DefaultDecoderConfig.Unmarshal([]byte(`1`), o1)
		checkError(t, EfaceLoopingError, err)
	})
	t.Run("loop 2", func(t *testing.T) {
		// the standard lib does not deal with cross nested
		type crossIface interface{}
		var o1 interface{}
		var o1o crossIface
		o1 = &o1o
		o1o = &o1
		err := DefaultDecoderConfig.Unmarshal([]byte(`1`), o1)
		checkError(t, EfaceLoopingError, err)
	})
}

func TestValDecoder_Native_Interface_Loop_WithDecoderConfig(t *testing.T) {
	t.Run("eface decoder", func(t *testing.T) {
		type crossIface interface{}
		var o1 interface{}
		var o1o crossIface
		o1 = &o1o
		o1o = &o1
		decCfg := NewDecoderConfig(&DecoderOption{
			ValDecoders: map[reflect.Type]ValDecoder{
				reflect.TypeOf((*interface{})(nil)).Elem(): (*testIfaceDecoder)(nil),
			},
		})
		err := decCfg.Unmarshal([]byte(`"abc"`), o1)
		require.NoError(t, err)
		printValue(t, ">>", o1)
		require.Equal(t, "abc", o1)
	})
	t.Run("eface decoder 2", func(t *testing.T) {
		type crossIface interface{}
		var o1 interface{}
		var o1o crossIface
		o1 = &o1o
		o1o = &o1
		decCfg := NewDecoderConfig(&DecoderOption{
			ValDecoders: map[reflect.Type]ValDecoder{
				reflect.TypeOf((*crossIface)(nil)).Elem(): (*testIfaceDecoder)(nil),
			},
		})
		err := decCfg.Unmarshal([]byte(`"abc"`), o1)
		require.NoError(t, err)
		printValue(t, ">>", o1)
		require.Equal(t, "abc", o1o)
	})
	debug.FreeOSMemory()
}
