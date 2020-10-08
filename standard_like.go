package jzon

import "io"

// UnmarshalFromReader behave like json.Unmarshal but with an io.Reader
func UnmarshalFromReader(r io.Reader, o interface{}) error {
	return DefaultDecoderConfig.UnmarshalFromReader(r, o)
}

// UnmarshalFromString behave like json.Unmarshal but with a string
func UnmarshalFromString(s string, o interface{}) error {
	return DefaultDecoderConfig.UnmarshalFromString(s, o)
}
