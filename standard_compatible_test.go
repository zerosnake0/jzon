package jzon

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

var (
	compatibleOnError = os.Getenv("COMPATIBLE_ON_ERROR") == "1"
)

func nestedArray1(count int) string {
	return strings.Repeat(" [", count) + " [ ] " +
		strings.Repeat("] ", count)
}

func nestedArray2(count int) string {
	return strings.Repeat(" [ [ ], ", count) + " [ ] " +
		strings.Repeat("] ", count)
}

func nestedArrayWithObject(count int) string {
	return strings.Repeat(" [ { }, ", count) + " [ ] " +
		strings.Repeat("] ", count)
}

func nestedObject(count int) string {
	return strings.Repeat(` { "a" : { }, "b": `, count) + " { } " +
		strings.Repeat("} ", count)
}

func nestedObjectWithArray(count int) string {
	return strings.Repeat(` { "a" : [ ], "b": `, count) + " [ ] " +
		strings.Repeat("} ", count)
}

type printValueKey struct {
	rtype rtype
	ptr   uintptr
}

func (pk printValueKey) String() string {
	return fmt.Sprintf("<%x %x>", pk.rtype, pk.ptr)
}

func printValue(t *testing.T, prefix string, o interface{}) {
	prefix += " "
	if o == nil {
		t.Logf(prefix + "nil")
		return
	}
	visited := map[printValueKey]bool{}
	oV := reflect.ValueOf(o)
	for indent := prefix; ; indent += "  " {
		i := oV.Interface()
		ef := (*eface)(unsafe.Pointer(&i))
		vk := printValueKey{ef.rtype, uintptr(ef.data)}

		k := oV.Kind()
		t.Logf(indent+"%+v %+v %v", oV.Type(), oV, vk)

		if k != reflect.Interface && k != reflect.Ptr {
			break
		}
		if oV.IsNil() {
			break
		}

		if visited[vk] {
			t.Logf(indent + "  visited...")
			break
		}
		visited[vk] = true

		oV = oV.Elem()
	}
}

func checkDecodeWithStandard(t *testing.T, decCfg *DecoderConfig, data string, ex error, exp, got interface{}) {
	b := []byte(data)
	expErr := json.Unmarshal(b, exp)
	gotErr := decCfg.Unmarshal(b, got)
	t.Logf("\nexpErr: %+v\ngotErr: %+v", expErr, gotErr)
	noError := expErr == nil
	if noError {
		printValue(t, "exp", reflect.ValueOf(exp).Elem().Interface())
	}
	require.Equal(t, noError, gotErr == nil,
		"exp %+v\ngot %+v", expErr, gotErr)
	require.Equalf(t, noError, ex == nil, "exp err: %v\ngot err: %v", ex, gotErr)
	if ex != nil {
		checkError(t, ex, gotErr)
		// if reflect.TypeOf(errors.New("")) == reflect.TypeOf(ex) {
		// 	require.Equalf(t, ex, gotErr, "exp err:%v\ngot err:%v", ex, gotErr)
		// } else {
		// 	require.IsTypef(t, ex, gotErr, "exp err:%v\ngot err:%v", ex, gotErr)
		// }
	}
	if !noError && !compatibleOnError {
		return
	}
	if exp == nil {
		require.Equal(t, nil, got)
		return
	}
	expV := reflect.ValueOf(exp)
	gotV := reflect.ValueOf(got)
	if expV.IsNil() {
		require.True(t, gotV.IsNil())
		return
	}
	expI := expV.Elem().Interface()
	gotI := gotV.Elem().Interface()
	printValue(t, "got", gotI)
	require.Equalf(t, expI, gotI, "exp %+v\ngot %+v", expI, gotI)
}

func TestValid(t *testing.T) {
	f := func(t *testing.T, s string) {
		data := localStringToBytes(s)
		require.Equal(t, json.Valid(data), Valid(data))
	}
	t.Run("empty", func(t *testing.T) {
		f(t, "")
	})
	t.Run("empty object", func(t *testing.T) {
		f(t, "{}")
	})
	t.Run("data remained", func(t *testing.T) {
		f(t, "{}1")
	})
}
