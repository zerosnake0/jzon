package jzon

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalFromReader(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		r := strings.NewReader("{}")
		var i interface{}
		err := UnmarshalFromReader(r, &i)
		require.NoError(t, err)
	})
	t.Run("failure", func(t *testing.T) {
		r := iotest.TimeoutReader(iotest.OneByteReader(strings.NewReader("{}")))
		var i interface{}
		err := UnmarshalFromReader(r, &i)
		checkError(t, iotest.ErrTimeout, err)
	})
	t.Run("final error", func(t *testing.T) {
		e := errors.New("test")
		r := &stepByteReader{
			b:   "{}",
			err: e,
		}
		var i interface{}
		err := UnmarshalFromReader(r, &i)
		checkError(t, e, err)
	})
}

func TestUnmarshalFromString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := "{}"
		var i interface{}
		err := UnmarshalFromString(s, &i)
		require.NoError(t, err)
	})
	t.Run("failure", func(t *testing.T) {
		s := "{"
		var i interface{}
		err := UnmarshalFromString(s, &i)
		checkError(t, io.EOF, err)
	})
}
