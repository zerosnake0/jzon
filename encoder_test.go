package jzon

import "encoding/json"

type encFace interface {
	Encode(v interface{}) error

	SetEscapeHTML(on bool)

	SetIndent(prefix, indent string)
}

var _ encFace = &json.Encoder{}
var _ encFace = &Encoder{}
