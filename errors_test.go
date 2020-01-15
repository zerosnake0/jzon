package jzon

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

var errorsNewType = reflect.TypeOf(errors.New(""))

func checkError(t *testing.T, exp, got error) {
	internalErr, ok := got.(*DecodeError)
	if ok {
		got = internalErr.reason
	}
	if reflect.TypeOf(exp) == errorsNewType {
		require.Equalf(t, exp, got, "exp err:%v\ngot err:%v", exp, got)
	} else {
		require.IsTypef(t, exp, got, "exp err:%v\ngot err:%v", exp, got)
	}
}
