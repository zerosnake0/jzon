package jzon

type Encoder struct {
}

func (e Encoder) Encode(v interface{}) error {
	panic("implement me")
}

func (e Encoder) SetEscapeHTML(on bool) {
	panic("implement me")
}

func (e Encoder) SetIndent(prefix, indent string) {
	panic("implement me")
}
