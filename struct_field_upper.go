package jzon

import (
	"unicode"
	"unicode/utf8"
)

func toUpper(in, out []byte) []byte {
	l := len(in)
	for i := 0; i < l; {
		r, size := utf8.DecodeRune(in[i:])
		fr := unicode.SimpleFold(r)
		for fr > r {
			fr = unicode.SimpleFold(fr)
		}
		out = appendRune(out, fr)
		i += size
	}
	return out
}
