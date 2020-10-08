package jzon

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/require"
)

type decFace interface {
	UseNumber()

	DisallowUnknownFields()

	Decode(interface{}) error

	Buffered() io.Reader

	// This is incompatible
	// Token() (json.Token, error)

	More() bool
}

type inputOffset interface {
	// Not available for go<1.13
	InputOffset() int64
}

var _ decFace = &json.Decoder{}
var _ decFace = &Decoder{}

func TestDecoder(t *testing.T) {
	t.Run("consecutive", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			must := require.New(t)

			newReader := func() io.Reader {
				s := ` {} {} `
				return iotest.OneByteReader(strings.NewReader(s))
			}

			check := func(dec decFace, length int, expMore bool, leftOffset, rightOffset int64) {
				buffered := dec.Buffered()
				if buffered == nil {
					return
				}
				b, err := ioutil.ReadAll(buffered)
				must.NoError(err)
				if length > 0 {
					t.Logf("%T: %q", dec, b)
				}
				must.True(length >= len(b))

				ofI, ofIok := dec.(inputOffset)
				if ofIok {
					offset := ofI.InputOffset()
					// t.Logf("%T %d", dec, offset)
					must.True(leftOffset <= offset, "%T %d > %d", dec, leftOffset, offset)
					must.True(rightOffset >= offset, "%T %d < %d", dec, rightOffset, offset)
				}

				more := dec.More()
				must.Equal(expMore, more, "%T", dec)

				if ofIok {
					offset := ofI.InputOffset()
					// t.Logf("%T %d", dec, offset)
					must.True(leftOffset <= offset, "%T %d > %d", dec, leftOffset, offset)
					must.True(rightOffset >= offset, "%T %d < %d", dec, rightOffset, offset)
				}
			}

			f := func(dec decFace) {
				var i, i2, i3 interface{}

				check(dec, 0, true, 0, 1)

				err := dec.Decode(&i)
				must.NoError(err, "%T", dec)
				check(dec, 0, true, 3, 4)

				err2 := dec.Decode(&i2)
				must.NoError(err2, "%T", dec)
				check(dec, 0, false, 6, 7)

				err3 := dec.Decode(&i3)
				must.Equal(io.EOF, err3, "%T", dec)
				check(dec, 1, false, 6, 7)
			}
			f(json.NewDecoder(newReader()))
			f(NewDecoder(newReader()))
		})
		t.Run("failure at start", func(t *testing.T) {
			must := require.New(t)

			newReader := func() io.Reader {
				s := ` } {} `
				return iotest.OneByteReader(strings.NewReader(s))
			}
			f := func(dec decFace) {
				var i, i2 interface{}

				err := dec.Decode(&i)
				t.Logf("%T %v", dec, err)
				must.Error(err)

				err2 := dec.Decode(&i2)
				t.Logf("%T %v", dec, err2)
				must.Equal(err2, err)
			}
			f(json.NewDecoder(newReader()))
			f(NewDecoder(newReader()))
		})
		t.Run("failure at middle", func(t *testing.T) {
			must := require.New(t)

			newReader := func() io.Reader {
				s := ` {} } `
				return iotest.OneByteReader(strings.NewReader(s))
			}
			f := func(dec decFace) {
				var i, i2 interface{}

				err := dec.Decode(&i)
				must.NoError(err)

				err2 := dec.Decode(&i2)
				t.Logf("%T %v", dec, err2)
				must.Error(err2)
			}
			f(json.NewDecoder(newReader()))
			f(NewDecoder(newReader()))
		})
	})
}

func TestDecoder_UseNumber(t *testing.T) {
	must := require.New(t)
	s := "123.456"
	newReader := func() io.Reader {
		return strings.NewReader(s)
	}
	f := func(dec decFace) {
		dec.UseNumber()
		var i interface{}
		err := dec.Decode(&i)
		must.NoError(err)
		must.Equal(json.Number(s), i)
	}
	f(json.NewDecoder(newReader()))
	f(NewDecoder(newReader()))
}

func TestDecoder_DisallowUnknownFields(t *testing.T) {
	must := require.New(t)
	s := `{"k":"v"}`
	newReader := func() io.Reader {
		return strings.NewReader(s)
	}
	f := func(dec decFace) {
		dec.DisallowUnknownFields()
		var i struct{}
		err := dec.Decode(&i)
		t.Log(err)
		must.Error(err)
	}
	f(json.NewDecoder(newReader()))
	f(NewDecoder(newReader()))
}
