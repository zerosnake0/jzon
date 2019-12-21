package jzon

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func printValue(t *testing.T, prefix string, o interface{}) {
	prefix += " "
	if o == nil {
		t.Logf(prefix + "nil")
		return
	}
	oV := reflect.ValueOf(o)
	for indent := prefix; ; indent += "  " {
		k := oV.Kind()
		t.Logf(indent+"%+v %+v", oV.Type(), oV)
		if k != reflect.Interface && k != reflect.Ptr {
			break
		}
		if oV.IsNil() {
			break
		}
		oV = oV.Elem()
	}
}

func checkStandard(t *testing.T, decoder *Decoder, data string, ex error, exp, got interface{}) {
	b := []byte(data)
	expErr := json.Unmarshal(b, exp)
	gotErr := decoder.Unmarshal(b, got)
	t.Logf("\nexpErr: %+v\ngotErr: %+v", expErr, gotErr)
	noError := expErr == nil
	if noError {
		printValue(t, "exp", reflect.ValueOf(exp).Elem().Interface())
	}
	require.Equal(t, noError, gotErr == nil,
		"exp %+v\ngot %+v", expErr, gotErr)
	require.Equalf(t, noError, ex == nil, "exp err: %v\ngot err: %v", ex, gotErr)
	if ex != nil {
		if assert.ObjectsAreEqual(reflect.TypeOf(errors.New("")), reflect.TypeOf(ex)) {
			require.Equalf(t, ex, gotErr, "exp err:%v\ngot err:%v", ex, gotErr)
		} else {
			require.IsTypef(t, ex, gotErr, "exp err:%v\ngot err:%v", ex, gotErr)
		}
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
