package jzon

import "io"

func UnmarshalFromReader(r io.Reader, o interface{}) error {
	return DefaultDecoderConfig.UnmarshalFromReader(r, o)
}

func UnmarshalFromString(s string, o interface{}) error {
	return DefaultDecoderConfig.UnmarshalFromString(s, o)
}
