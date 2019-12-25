package jzon

import (
	"io"
)

type oneByteReader struct {
	b   []byte
	err error
}

var _ io.Reader = &oneByteReader{}

func (o *oneByteReader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	if len(o.b) == 0 {
		if o.err != nil {
			return 0, o.err
		}
		return 0, io.EOF
	}
	p[0] = o.b[0]
	o.b = o.b[1:]
	return 1, nil
}

type repeatByteReader struct {
	b     byte
	count int
}

var _ io.Reader = &repeatByteReader{}

func (r *repeatByteReader) Read(p []byte) (n int, err error) {
	l := len(p)
	for r.count > 0 && n < l {
		p[n] = r.b
		r.count--
		n++
	}
	return
}
