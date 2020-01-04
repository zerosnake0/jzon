package jzon

import (
	"encoding/json"
	"errors"
	"testing"
)

type testJsonMarshaler struct {
	data string
	err  error
}

func (m testJsonMarshaler) MarshalJSON() ([]byte, error) {
	return []byte(m.data), m.err
}

type testJsonMarshaler2 struct {
	data string
	err  error
}

func (m *testJsonMarshaler2) MarshalJSON() ([]byte, error) {
	return []byte(m.data), m.err
}

func TestValEncoder_JsonMarshaler_Error(t *testing.T) {
	t.Run("chain error", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			jsonMarshalerEncoder(0).Encode(nil, s, nil)
		})
	})
	t.Run("chain error (dynamic)", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			(*dynamicJsonMarshalerEncoder)(nil).Encode(nil, s, nil)
		})
	})
}

func TestValEncoder_JsonMarshaler_NonPointerReceiver(t *testing.T) {
	f := checkEncodeValueWithStandard
	t.Run("non pointer", func(t *testing.T) {
		t.Run("no error", func(t *testing.T) {
			f(t, testJsonMarshaler{
				data: `{"a":1}`,
			}, nil)
		})
		t.Run("error", func(t *testing.T) {
			e := errors.New("test")
			f(t, testJsonMarshaler{
				data: `{"a":1}`,
				err:  e,
			}, e)
		})
	})
	t.Run("pointer", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			f(t, (*testJsonMarshaler)(nil), nil)
		})
		t.Run("no error", func(t *testing.T) {
			f(t, &testJsonMarshaler{
				data: `{"a":2}`,
			}, nil)
		})
		t.Run("error", func(t *testing.T) {
			e := errors.New("test")
			f(t, &testJsonMarshaler{
				data: `{"a":2}`,
				err:  e,
			}, e)
		})
	})
	t.Run("pointer of pointer", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			f(t, (**testJsonMarshaler)(nil), nil)
		})
		t.Run("pointer of nil", func(t *testing.T) {
			ptr := (*testJsonMarshaler)(nil)
			f(t, &ptr, nil)
		})
		t.Run("no error", func(t *testing.T) {
			ptr := &testJsonMarshaler{
				data: `{"a":2}`,
			}
			f(t, &ptr, nil)
		})
		t.Run("error", func(t *testing.T) {
			e := errors.New("test")
			ptr := &testJsonMarshaler{
				data: `{"a":2}`,
				err:  e,
			}
			f(t, &ptr, e)
		})
	})
}

func TestValEncoder_JsonMarshaler_PointerReceiver(t *testing.T) {
	f := checkEncodeValueWithStandard
	t.Run("non pointer", func(t *testing.T) {
		t.Run("no error", func(t *testing.T) {
			f(t, testJsonMarshaler2{
				data: `{"b":1}`,
			}, nil)
		})
		t.Run("error", func(t *testing.T) {
			e := errors.New("test")
			f(t, testJsonMarshaler2{
				data: `{"b":1}`,
				err:  e,
			}, nil)
		})
	})
	t.Run("pointer", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			f(t, (*testJsonMarshaler2)(nil), nil)
		})
		t.Run("no error", func(t *testing.T) {
			f(t, &testJsonMarshaler2{
				data: `{"b":1}`,
			}, nil)
		})
		t.Run("error", func(t *testing.T) {
			e := errors.New("test")
			f(t, &testJsonMarshaler2{
				data: `{"b":1}`,
				err:  e,
			}, e)
		})
	})
	t.Run("pointer of pointer", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			f(t, (**testJsonMarshaler2)(nil), nil)
		})
		t.Run("pointer of nil", func(t *testing.T) {
			ptr := (*testJsonMarshaler2)(nil)
			f(t, &ptr, nil)
		})
		t.Run("no error", func(t *testing.T) {
			ptr := &testJsonMarshaler2{
				data: `{"a":2}`,
			}
			f(t, &ptr, nil)
		})
		t.Run("error", func(t *testing.T) {
			e := errors.New("test")
			ptr := &testJsonMarshaler2{
				data: `{"a":2}`,
				err:  e,
			}
			f(t, &ptr, e)
		})
	})
	t.Run("struct member", func(t *testing.T) {
		type st struct {
			A testJsonMarshaler2
		}
		/*
		 * with the current implementation,
		 * only one of the following two test can succeed
		 */
		t.Run("value", func(t *testing.T) {
			skipTest(t, "pointer encoder on value")
			checkEncodeValueWithStandard(t, st{
				A: testJsonMarshaler2{
					data: `{"a":2}`,
				},
			}, nil)
		})
		t.Run("ptr", func(t *testing.T) {
			checkEncodeValueWithStandard(t, &st{
				A: testJsonMarshaler2{
					data: `{"a":2}`,
				},
			}, nil)
		})
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
			data: `"test"`,
			err:  e,
		}
		checkEncodeValueWithStandard(t, &i, e)
	})
	t.Run("marshaler", func(t *testing.T) {
		var i json.Marshaler = testJsonMarshaler{
			data: `"test"`,
		}
		checkEncodeValueWithStandard(t, &i, nil)
	})
	t.Run("marshaler 2", func(t *testing.T) {
		var i json.Marshaler = &testJsonMarshaler{
			data: `"test 2"`,
		}
		checkEncodeValueWithStandard(t, &i, nil)
	})
}
