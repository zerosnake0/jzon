package jzon

import (
	"encoding"
	"errors"
	"testing"
)

type testTextMarshaler struct {
	data []byte
	err  error
}

func (m testTextMarshaler) MarshalText() ([]byte, error) {
	return m.data, m.err
}

type testTextMarshaler2 struct {
	data []byte
	err  error
}

func (m *testTextMarshaler2) MarshalText() ([]byte, error) {
	return m.data, m.err
}

func TestValEncoder_TextMarshaler(t *testing.T) {
	f := func(t *testing.T, m encoding.TextMarshaler) {
		checkEncodeValueWithStandard(t, DefaultEncoder, m)
	}
	t.Run("non pointer receiver", func(t *testing.T) {
		t.Run("non pointer", func(t *testing.T) {
			f(t, testTextMarshaler{
				data: []byte(`{"a":1}`),
			})
		})
		t.Run("non pointer error", func(t *testing.T) {
			f(t, testTextMarshaler{
				data: []byte(`{"a":1}`),
				err:  errors.New("test"),
			})
		})
		t.Run("pointer", func(t *testing.T) {
			f(t, &testTextMarshaler{
				data: []byte(`{"a":2}`),
			})
		})
		t.Run("pointer error", func(t *testing.T) {
			f(t, &testTextMarshaler{
				data: []byte(`{"a":2}`),
				err:  errors.New("test"),
			})
		})
		t.Run("nil pointer", func(t *testing.T) {
			f(t, (*testTextMarshaler)(nil))
		})
	})
	t.Run("pointer receiver", func(t *testing.T) {
		t.Run("pointer", func(t *testing.T) {
			f(t, &testTextMarshaler2{
				data: []byte(`{"b":1}`),
			})
		})
		t.Run("pointer error", func(t *testing.T) {
			f(t, &testTextMarshaler2{
				data: []byte(`{"b":1}`),
				err:  errors.New("test"),
			})
		})
		t.Run("nil pointer", func(t *testing.T) {
			f(t, (*testTextMarshaler2)(nil))
		})
	})
}

func TestValEncoder_DynamicTextMarshaler(t *testing.T) {
	t.Run("marshaler <nil>", func(t *testing.T) {
		// TODO: This test should be automatically fixed in the future golang version
		v := "go1.13.5"
		if goVersion.LessEqual(v) {
			t.Skipf("skipping this test for go version <= %s", v)
		}
		var i encoding.TextMarshaler
		checkEncodeValueWithStandard(t, DefaultEncoder, &i)
	})
	t.Run("marshaler error", func(t *testing.T) {
		var i encoding.TextMarshaler = testTextMarshaler{
			data: []byte(`"test"`),
			err:  errors.New("test"),
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, &i)
	})
	t.Run("marshaler", func(t *testing.T) {
		var i encoding.TextMarshaler = testTextMarshaler{
			data: []byte(`"test"`),
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, &i)
	})
	t.Run("marshaler 2", func(t *testing.T) {
		var i encoding.TextMarshaler = &testTextMarshaler{
			data: []byte(`"test 2"`),
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, &i)
	})
}
