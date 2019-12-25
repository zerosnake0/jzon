package jzon

import (
	"encoding/json"
	"errors"
	"testing"
)

type testJsonMarshaler struct {
	data []byte
	err  error
}

func (m testJsonMarshaler) MarshalJSON() ([]byte, error) {
	return m.data, m.err
}

type testJsonMarshaler2 struct {
	data []byte
	err  error
}

func (m *testJsonMarshaler2) MarshalJSON() ([]byte, error) {
	return m.data, m.err
}

func TestValEncoder_JsonMarshaler_NonPointerReceiver(t *testing.T) {
	f := func(t *testing.T, m json.Marshaler) {
		checkEncodeValueWithStandard(t, DefaultEncoder, m)
	}
	t.Run("non pointer", func(t *testing.T) {
		f(t, testJsonMarshaler{
			data: []byte(`{"a":1}`),
		})
	})
	t.Run("non pointer error", func(t *testing.T) {
		f(t, testJsonMarshaler{
			data: []byte(`{"a":1}`),
			err:  errors.New("test"),
		})
	})
	t.Run("pointer", func(t *testing.T) {
		f(t, &testJsonMarshaler{
			data: []byte(`{"a":2}`),
		})
	})
	t.Run("pointer error", func(t *testing.T) {
		f(t, &testJsonMarshaler{
			data: []byte(`{"a":2}`),
			err:  errors.New("test"),
		})
	})
	t.Run("nil pointer", func(t *testing.T) {
		f(t, (*testJsonMarshaler)(nil))
	})
}

func TestValEncoder_JsonMarshaler_PointerReceiver(t *testing.T) {
	f := func(t *testing.T, m json.Marshaler) {
		checkEncodeValueWithStandard(t, DefaultEncoder, m)
	}
	t.Run("pointer", func(t *testing.T) {
		f(t, &testJsonMarshaler2{
			data: []byte(`{"b":1}`),
		})
	})
	t.Run("pointer error", func(t *testing.T) {
		f(t, &testJsonMarshaler2{
			data: []byte(`{"b":1}`),
			err:  errors.New("test"),
		})
	})
	t.Run("nil pointer", func(t *testing.T) {
		f(t, (*testJsonMarshaler2)(nil))
	})
}

func TestValEncoder_DynamicJsonMarshaler(t *testing.T) {
	t.Run("marshaler <nil>", func(t *testing.T) {
		var i json.Marshaler
		checkEncodeValueWithStandard(t, DefaultEncoder, &i)
	})
	t.Run("marshaler error", func(t *testing.T) {
		var i json.Marshaler = testJsonMarshaler{
			data: []byte(`"test"`),
			err:  errors.New("test"),
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, &i)
	})
	t.Run("marshaler", func(t *testing.T) {
		var i json.Marshaler = testJsonMarshaler{
			data: []byte(`"test"`),
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, &i)
	})
	t.Run("marshaler 2", func(t *testing.T) {
		var i json.Marshaler = &testJsonMarshaler{
			data: []byte(`"test 2"`),
		}
		checkEncodeValueWithStandard(t, DefaultEncoder, &i)
	})
}
