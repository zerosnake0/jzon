package jzon

import (
	"encoding/json"
	"io"
)

type decFace interface {
	UseNumber()

	DisallowUnknownFields()

	Decode(interface{}) error

	Buffered() io.Reader

	// This is incompatible
	// Token() (json.Token, error)

	More() bool

	InputOffset() int64
}

var _ decFace = &json.Decoder{}
var _ decFace = &Decoder{}
