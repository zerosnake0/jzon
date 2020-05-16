package jzon

import (
	"encoding"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValEncoder_TextMarshaler_ChainError(t *testing.T) {
	t.Run("direct", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			(*textMarshalerEncoder)(nil).Encode(nil, s, nil)
		})
	})
	t.Run("dynamic", func(t *testing.T) {
		testStreamerChainError(t, func(s *Streamer) {
			(*dynamicTextMarshalerEncoder)(nil).Encode(nil, s, nil)
		})
	})
}

type testTextMarshaler struct {
	data string
	err  error
}

func (m testTextMarshaler) MarshalText() ([]byte, error) {
	return []byte(m.data), m.err
}

type testTextMarshaler2 struct {
	data string
	err  error
}

func (m *testTextMarshaler2) MarshalText() ([]byte, error) {
	if m == nil {
		return []byte(`is_null`), nil
	}
	return []byte(m.data), m.err
}

func TestValEncoder_TextMarshaler(t *testing.T) {
	f := func(t *testing.T, m encoding.TextMarshaler, err error) {
		checkEncodeValueWithStandard(t, m, err)
	}
	t.Run("non pointer receiver", func(t *testing.T) {
		t.Run("non pointer", func(t *testing.T) {
			f(t, testTextMarshaler{
				data: `{"a":1}`,
			}, nil)
		})
		t.Run("non pointer error", func(t *testing.T) {
			e := errors.New("test")
			f(t, testTextMarshaler{
				data: `{"a":1}`,
				err:  e,
			}, e)
		})
		t.Run("pointer", func(t *testing.T) {
			f(t, &testTextMarshaler{
				data: `{"a":2}`,
			}, nil)
		})
		t.Run("pointer error", func(t *testing.T) {
			e := errors.New("test")
			f(t, &testTextMarshaler{
				data: `{"a":2}`,
				err:  e,
			}, e)
		})
		t.Run("nil pointer", func(t *testing.T) {
			f(t, (*testTextMarshaler)(nil), nil)
		})
	})
	t.Run("pointer receiver", func(t *testing.T) {
		t.Run("pointer", func(t *testing.T) {
			f(t, &testTextMarshaler2{
				data: `{"b":1}`,
			}, nil)
		})
		t.Run("pointer error", func(t *testing.T) {
			e := errors.New("test")
			f(t, &testTextMarshaler2{
				data: `{"b":1}`,
				err:  e,
			}, e)
		})
		t.Run("nil pointer", func(t *testing.T) {
			f(t, (*testTextMarshaler2)(nil), nil)
		})
	})
}

func TestValEncoder_DynamicTextMarshaler(t *testing.T) {
	t.Run("marshaler nil", func(t *testing.T) {
		// TODO: This test should be automatically fixed in the future golang version
		v := "go1.13.11"
		if goVersion.LessEqual(v) {
			var i encoding.TextMarshaler
			b, err := Marshal(&i)
			require.NoError(t, err)
			require.Equal(t, "null", string(b))
		} else {
			var i encoding.TextMarshaler
			checkEncodeValueWithStandard(t, &i, nil)
		}
	})
	t.Run("marshaler error", func(t *testing.T) {
		e := errors.New("test")
		var i encoding.TextMarshaler = testTextMarshaler{
			data: `"test"`,
			err:  e,
		}
		checkEncodeValueWithStandard(t, &i, e)
	})
	t.Run("marshaler", func(t *testing.T) {
		var i encoding.TextMarshaler = testTextMarshaler{
			data: `"test"`,
		}
		checkEncodeValueWithStandard(t, &i, nil)
	})
	t.Run("marshaler 2", func(t *testing.T) {
		var i encoding.TextMarshaler = &testTextMarshaler{
			data: `"test 2"`,
		}
		checkEncodeValueWithStandard(t, &i, nil)
	})
}

type testDirectTextMarshaler map[int]int

func (m testDirectTextMarshaler) MarshalText() ([]byte, error) {
	s := fmt.Sprintf("%d", len(m))
	return []byte(s), nil
}

func TestValEncoder_TextMarshaler_Direct(t *testing.T) {
	t.Run("value", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			checkEncodeValueWithStandard(t, testDirectTextMarshaler(nil), nil)
		})
		t.Run("non nil", func(t *testing.T) {
			checkEncodeValueWithStandard(t, testDirectTextMarshaler{
				1: 2,
			}, nil)
		})
	})
	t.Run("pointer", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			checkEncodeValueWithStandard(t, (*testDirectTextMarshaler)(nil), nil)
		})
		t.Run("non nil", func(t *testing.T) {
			var m testDirectTextMarshaler
			checkEncodeValueWithStandard(t, &m, nil)
		})
		t.Run("non nil 2", func(t *testing.T) {
			checkEncodeValueWithStandard(t, &testDirectTextMarshaler{
				1: 2,
			}, nil)
		})
	})
	t.Run("struct member", func(t *testing.T) {
		type st struct {
			A testDirectTextMarshaler
		}
		checkEncodeValueWithStandard(t, &st{}, nil)
	})
}

func TestValEncoder_TextMarshaler_OmitEmpty(t *testing.T) {
	t.Run("text marshaler", func(t *testing.T) {
		type st struct {
			A testTextMarshaler `json:",omitempty"`
		}
		checkEncodeValueWithStandard(t, st{
			A: testTextMarshaler{
				data: "true",
			},
		}, nil)
	})
	t.Run("direct text marshaler", func(t *testing.T) {
		type st struct {
			A testDirectTextMarshaler `json:",omitempty"`
		}
		t.Run("nil", func(t *testing.T) {
			checkEncodeValueWithStandard(t, st{}, nil)
		})
		t.Run("zero", func(t *testing.T) {
			checkEncodeValueWithStandard(t, st{
				A: testDirectTextMarshaler{},
			}, nil)
		})
		t.Run("non zero", func(t *testing.T) {
			checkEncodeValueWithStandard(t, st{
				A: testDirectTextMarshaler{1: 2},
			}, nil)
		})
	})
	t.Run("pointer text marshaler", func(t *testing.T) {
		type st struct {
			A testTextMarshaler2 `json:",omitempty"`
		}
		checkEncodeValueWithStandard(t, &st{
			A: testTextMarshaler2{
				data: "true",
			},
		}, nil)
	})
	t.Run("dynamic text marshaler", func(t *testing.T) {
		type st struct {
			A encoding.TextMarshaler `json:",omitempty"`
		}
		t.Run("nil", func(t *testing.T) {
			checkEncodeValueWithStandard(t, st{}, nil)
		})
		t.Run("nil pointer", func(t *testing.T) {
			checkEncodeValueWithStandard(t, st{
				A: (*testTextMarshaler2)(nil),
			}, nil)
		})
	})
}
