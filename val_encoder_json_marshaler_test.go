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

func TestValEncoder_JsonMarshaler_Error(t *testing.T) {
	t.Run("chain error", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			jsonMarshalerEncoder(0).Encode(nil, s)
		})
	})
	t.Run("chain error (dynamic)", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			(*dynamicJsonMarshalerEncoder)(nil).Encode(nil, s)
		})
	})
}

func TestValEncoder_JsonMarshaler_NonPointerReceiver(t *testing.T) {
	f := func(t *testing.T, m json.Marshaler, err error) {
		checkEncodeValueWithStandard(t, m, err)
	}
	t.Run("non pointer", func(t *testing.T) {
		f(t, testJsonMarshaler{
			data: []byte(`{"a":1}`),
		}, nil)
	})
	t.Run("non pointer error", func(t *testing.T) {
		e := errors.New("test")
		f(t, testJsonMarshaler{
			data: []byte(`{"a":1}`),
			err:  e,
		}, e)
	})
	t.Run("pointer", func(t *testing.T) {
		f(t, &testJsonMarshaler{
			data: []byte(`{"a":2}`),
		}, nil)
	})
	t.Run("pointer error", func(t *testing.T) {
		e := errors.New("test")
		f(t, &testJsonMarshaler{
			data: []byte(`{"a":2}`),
			err:  e,
		}, e)
	})
	t.Run("nil pointer", func(t *testing.T) {
		f(t, (*testJsonMarshaler)(nil), nil)
	})
}

func TestValEncoder_JsonMarshaler_PointerReceiver(t *testing.T) {
	f := func(t *testing.T, m json.Marshaler, err error) {
		checkEncodeValueWithStandard(t, m, err)
	}
	t.Run("pointer", func(t *testing.T) {
		f(t, &testJsonMarshaler2{
			data: []byte(`{"b":1}`),
		}, nil)
	})
	t.Run("pointer error", func(t *testing.T) {
		e := errors.New("test")
		f(t, &testJsonMarshaler2{
			data: []byte(`{"b":1}`),
			err:  e,
		}, e)
	})
	t.Run("nil pointer", func(t *testing.T) {
		f(t, (*testJsonMarshaler2)(nil), nil)
	})
}

func TestValEncoder_DynamicJsonMarshaler(t *testing.T) {
	t.Run("marshaler (nil)", func(t *testing.T) {
		var i json.Marshaler
		checkEncodeValueWithStandard(t, &i, nil)
	})
	t.Run("marshaler error", func(t *testing.T) {
		e := errors.New("test")
		var i json.Marshaler = testJsonMarshaler{
			data: []byte(`"test"`),
			err:  e,
		}
		checkEncodeValueWithStandard(t, &i, e)
	})
	t.Run("marshaler", func(t *testing.T) {
		var i json.Marshaler = testJsonMarshaler{
			data: []byte(`"test"`),
		}
		checkEncodeValueWithStandard(t, &i, nil)
	})
	t.Run("marshaler 2", func(t *testing.T) {
		var i json.Marshaler = &testJsonMarshaler{
			data: []byte(`"test 2"`),
		}
		checkEncodeValueWithStandard(t, &i, nil)
	})
}
