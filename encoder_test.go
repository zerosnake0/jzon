package jzon

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type encFace interface {
	Encode(v interface{}) error

	SetEscapeHTML(on bool)

	// This is incompatible with standard library
	// SetIndent(prefix, indent string)
}

var _ encFace = &json.Encoder{}
var _ encFace = &Encoder{}

func TestEncoder_SetEscapeHTML(t *testing.T) {
	f := func(t *testing.T, o interface{}, expOn, expOff string) {
		must := require.New(t)
		buf := bytes.NewBuffer(nil)
		f := func(enc encFace) {
			// enabled
			buf.Reset()
			err := enc.Encode(o)
			must.NoError(err)
			must.Equal(expOn+"\n", buf.String(), "%T", enc)

			// disabled
			buf.Reset()
			enc.SetEscapeHTML(false)
			err = enc.Encode(o)
			must.NoError(err)
			must.Equal(expOff+"\n", buf.String(), "%T", enc)
		}
		f(json.NewEncoder(buf))
		enc := NewEncoder(buf)
		defer enc.Release()
		f(enc)
	}
	t.Run("normal", func(t *testing.T) {
		s := "<>&"
		expOn := `"\u003c\u003e\u0026"`
		expOff := `"` + s + `"`
		f(t, s, expOn, expOff)
	})
	t.Run("struct", func(t *testing.T) {
		s := struct {
			A int `json:"&"`
		}{}
		expOn := `{"\u0026":0}`
		expOff := `{"&":0}`
		f(t, s, expOn, expOff)
	})
}

// func TestEncoder_SetIndent(t *testing.T) {
// 	must := require.New(t)
// 	s := map[string]interface{}{
// 		"k": "v",
// 	}
// 	buf := bytes.NewBuffer(nil)
// 	f := func(enc encFace) {
// 		// disabled
// 		buf.Reset()
// 		err := enc.Encode(s)
// 		must.NoError(err)
// 		must.Equal(`{"k":"v"}`+"\n", buf.String(), "%T", enc)
//
// 		// disabled
// 		buf.Reset()
// 		enc.SetIndent("p", "i")
// 		err = enc.Encode(s)
// 		must.NoError(err)
// 		t.Logf("\n%s", buf.Bytes())
// 		must.Equal("{\npi\"k\": \"v\"\np}\n", buf.String(), "%T", enc)
// 	}
// 	f(json.NewEncoder(buf))
// 	f(NewEncoder(buf))
// }
