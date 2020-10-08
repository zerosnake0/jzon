package jzon

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValEncoder_JsonMarshaler_Error(t *testing.T) {
	t.Run("chain error", func(t *testing.T) {
		t.Run("", func(t *testing.T) {
			testStreamerChainError(t, func(s *Streamer) {
				(*jsonMarshalerEncoder)(nil).Encode(nil, s, nil)
			})
		})
		t.Run("dynamic", func(t *testing.T) {
			testStreamerChainError(t, func(s *Streamer) {
				(*dynamicJSONMarshalerEncoder)(nil).Encode(nil, s, nil)
			})
		})
	})
}

func TestValEncoder_JsonMarshaler_NonPointerReceiver(t *testing.T) {
	f := checkEncodeValueWithStandard
	t.Run("non pointer", func(t *testing.T) {
		t.Run("no error", func(t *testing.T) {
			f(t, testJSONMarshaler{
				data: `{"a":1}`,
			}, nil)
		})
		t.Run("error", func(t *testing.T) {
			e := errors.New("test")
			f(t, testJSONMarshaler{
				data: `{"a":1}`,
				err:  e,
			}, e)
		})
	})
	t.Run("pointer", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			f(t, (*testJSONMarshaler)(nil), nil)
		})
		t.Run("no error", func(t *testing.T) {
			f(t, &testJSONMarshaler{
				data: `{"a":2}`,
			}, nil)
		})
		t.Run("error", func(t *testing.T) {
			e := errors.New("test")
			f(t, &testJSONMarshaler{
				data: `{"a":2}`,
				err:  e,
			}, e)
		})
	})
	t.Run("pointer of pointer", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			f(t, (**testJSONMarshaler)(nil), nil)
		})
		t.Run("pointer of nil", func(t *testing.T) {
			ptr := (*testJSONMarshaler)(nil)
			f(t, &ptr, nil)
		})
		t.Run("no error", func(t *testing.T) {
			ptr := &testJSONMarshaler{
				data: `{"a":2}`,
			}
			f(t, &ptr, nil)
		})
		t.Run("error", func(t *testing.T) {
			e := errors.New("test")
			ptr := &testJSONMarshaler{
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
			f(t, testJSONMarshaler2{
				data: `{"b":1}`,
			}, nil)
		})
		t.Run("error", func(t *testing.T) {
			e := errors.New("test")
			f(t, testJSONMarshaler2{
				data: `{"b":1}`,
				err:  e,
			}, nil)
		})
	})
	t.Run("pointer", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			f(t, (*testJSONMarshaler2)(nil), nil)
		})
		t.Run("no error", func(t *testing.T) {
			f(t, &testJSONMarshaler2{
				data: `{"b":1}`,
			}, nil)
		})
		t.Run("error", func(t *testing.T) {
			e := errors.New("test")
			f(t, &testJSONMarshaler2{
				data: `{"b":1}`,
				err:  e,
			}, e)
		})
	})
	t.Run("pointer of pointer", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			f(t, (**testJSONMarshaler2)(nil), nil)
		})
		t.Run("pointer of nil", func(t *testing.T) {
			ptr := (*testJSONMarshaler2)(nil)
			f(t, &ptr, nil)
		})
		t.Run("no error", func(t *testing.T) {
			ptr := &testJSONMarshaler2{
				data: `{"a":2}`,
			}
			f(t, &ptr, nil)
		})
		t.Run("error", func(t *testing.T) {
			e := errors.New("test")
			ptr := &testJSONMarshaler2{
				data: `{"a":2}`,
				err:  e,
			}
			f(t, &ptr, e)
		})
	})
	t.Run("struct member", func(t *testing.T) {
		t.Run("value", func(t *testing.T) {
			type st struct {
				A testJSONMarshaler2
			}
			/*
			 * with the current implementation,
			 * only one of the following two test can succeed
			 */
			t.Run("value", func(t *testing.T) {
				skipTest(t, "pointer encoder on value")
				checkEncodeValueWithStandard(t, st{
					A: testJSONMarshaler2{
						data: `{"a":2}`,
					},
				}, nil)
			})
			t.Run("ptr", func(t *testing.T) {
				checkEncodeValueWithStandard(t, &st{
					A: testJSONMarshaler2{
						data: `{"a":2}`,
					},
				}, nil)
			})
		})
		t.Run("pointer", func(t *testing.T) {
			type st struct {
				A *testJSONMarshaler2
			}
			t.Run("nil", func(t *testing.T) {
				checkEncodeValueWithStandard(t, &st{}, nil)
			})
		})
	})
}

func TestValEncoder_DynamicJsonMarshaler(t *testing.T) {
	t.Run("marshaler nil", func(t *testing.T) {
		var i json.Marshaler
		checkEncodeValueWithStandard(t, &i, nil)
	})
	t.Run("marshaler error", func(t *testing.T) {
		e := errors.New("test")
		var i json.Marshaler = testJSONMarshaler{
			data: `"test"`,
			err:  e,
		}
		checkEncodeValueWithStandard(t, &i, e)
	})
	t.Run("marshaler", func(t *testing.T) {
		var i json.Marshaler = testJSONMarshaler{
			data: `"test"`,
		}
		checkEncodeValueWithStandard(t, &i, nil)
	})
	t.Run("marshaler 2", func(t *testing.T) {
		var i json.Marshaler = &testJSONMarshaler{
			data: `"test 2"`,
		}
		checkEncodeValueWithStandard(t, &i, nil)
	})
}

func TestValEncoder_JsonMarshaler_Direct(t *testing.T) {
	t.Run("value", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			checkEncodeValueWithStandard(t, testMapJSONMarshaler(nil), nil)
		})
		t.Run("non nil", func(t *testing.T) {
			checkEncodeValueWithStandard(t, testMapJSONMarshaler{
				1: 2,
			}, nil)
		})
	})
	t.Run("pointer", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			checkEncodeValueWithStandard(t, (*testMapJSONMarshaler)(nil), nil)
		})
		t.Run("non nil", func(t *testing.T) {
			var m testMapJSONMarshaler
			checkEncodeValueWithStandard(t, &m, nil)
		})
		t.Run("non nil 2", func(t *testing.T) {
			checkEncodeValueWithStandard(t, &testMapJSONMarshaler{
				1: 2,
			}, nil)
		})
	})
	t.Run("struct member", func(t *testing.T) {
		type st struct {
			A testMapJSONMarshaler
		}
		checkEncodeValueWithStandard(t, &st{}, nil)
	})
	t.Run("value of marshaler", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			var m json.Marshaler
			checkEncodeValueWithStandard(t, m, nil)
		})
		t.Run("value", func(t *testing.T) {
			var m json.Marshaler = testMapJSONMarshaler{
				1: 2,
			}
			checkEncodeValueWithStandard(t, m, nil)
		})
		t.Run("pointer", func(t *testing.T) {
			var m json.Marshaler = &testMapJSONMarshaler{
				1: 2,
			}
			checkEncodeValueWithStandard(t, m, nil)
		})
	})
	t.Run("pointer of marshaler", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			var m json.Marshaler
			checkEncodeValueWithStandard(t, &m, nil)
		})
		t.Run("value", func(t *testing.T) {
			var m json.Marshaler = testMapJSONMarshaler{
				1: 2,
			}
			checkEncodeValueWithStandard(t, &m, nil)
		})
		t.Run("pointer", func(t *testing.T) {
			var m json.Marshaler = &testMapJSONMarshaler{
				1: 2,
			}
			checkEncodeValueWithStandard(t, &m, nil)
		})
	})
}

func TestValEncoder_JsonMarshaler_OmitEmpty(t *testing.T) {
	t.Run("json marshaler", func(t *testing.T) {
		type st struct {
			A testJSONMarshaler `json:",omitempty"`
		}
		checkEncodeValueWithStandard(t, st{
			A: testJSONMarshaler{
				data: "true",
			},
		}, nil)
	})
	t.Run("indirect json marshaler", func(t *testing.T) {
		t.Run("bool", func(t *testing.T) {
			type st struct {
				A testBoolJSONMarshaler `json:",omitempty"`
			}
			require.True(t, ifaceIndir(rtypeOfType(reflect.TypeOf(st{}).Field(0).Type)))
			checkEncodeValueWithStandard(t, st{
				A: true,
			}, nil)
			checkEncodeValueWithStandard(t, st{
				A: false,
			}, nil)
		})
		t.Run("array", func(t *testing.T) {
			type st struct {
				A testIndirectArrayMarshaler `json:",omitempty"`
			}
			require.True(t, ifaceIndir(rtypeOfType(reflect.TypeOf(st{}).Field(0).Type)))
			checkEncodeValueWithStandard(t, st{}, nil)
			checkEncodeValueWithStandard(t, st{
				A: testIndirectArrayMarshaler{2},
			}, nil)
		})
		t.Run("slice", func(t *testing.T) {
			type st struct {
				A testSliceMarshaler `json:",omitempty"`
			}
			require.True(t, ifaceIndir(rtypeOfType(reflect.TypeOf(st{}).Field(0).Type)))
			t.Run("nil", func(t *testing.T) {
				checkEncodeValueWithStandard(t, st{}, nil)
			})
			t.Run("empty", func(t *testing.T) {
				checkEncodeValueWithStandard(t, st{
					A: testSliceMarshaler{},
				}, nil)
			})
			t.Run("non empty", func(t *testing.T) {
				checkEncodeValueWithStandard(t, st{
					A: testSliceMarshaler{4, 5, 6},
				}, nil)
			})
		})
		t.Run("struct", func(t *testing.T) {
			type st struct {
				A testIndirectStructMarshaler `json:",omitempty"`
			}
			require.True(t, ifaceIndir(rtypeOfType(reflect.TypeOf(st{}).Field(0).Type)))
			checkEncodeValueWithStandard(t, st{}, nil)
			checkEncodeValueWithStandard(t, st{
				A: testIndirectStructMarshaler{1},
			}, nil)
		})
	})
	t.Run("direct json marshaler", func(t *testing.T) {
		t.Run("array", func(t *testing.T) {
			type st struct {
				A testDirectArrayMarshaler `json:",omitempty"`
			}
			require.False(t, ifaceIndir(rtypeOfType(reflect.TypeOf(st{}).Field(0).Type)))
			checkEncodeValueWithStandard(t, st{}, nil)
			i := 123
			checkEncodeValueWithStandard(t, st{
				A: testDirectArrayMarshaler{&i},
			}, nil)
		})
		t.Run("map", func(t *testing.T) {
			type st struct {
				A testMapJSONMarshaler `json:",omitempty"`
			}
			require.False(t, ifaceIndir(rtypeOfType(reflect.TypeOf(st{}).Field(0).Type)))
			t.Run("nil", func(t *testing.T) {
				checkEncodeValueWithStandard(t, st{}, nil)
			})
			t.Run("zero", func(t *testing.T) {
				checkEncodeValueWithStandard(t, st{
					A: testMapJSONMarshaler{},
				}, nil)
			})
			t.Run("non zero", func(t *testing.T) {
				checkEncodeValueWithStandard(t, st{
					A: testMapJSONMarshaler{1: 2},
				}, nil)
			})
		})
		t.Run("struct", func(t *testing.T) {
			type st struct {
				A testDirectStructMarshaler `json:",omitempty"`
			}
			require.False(t, ifaceIndir(rtypeOfType(reflect.TypeOf(st{}).Field(0).Type)))
			checkEncodeValueWithStandard(t, st{}, nil)
			i := 2
			checkEncodeValueWithStandard(t, st{
				A: testDirectStructMarshaler{&i},
			}, nil)
		})
	})
	t.Run("pointer json marshaler", func(t *testing.T) {
		t.Run("value field", func(t *testing.T) {
			type st struct {
				A testJSONMarshaler2 `json:",omitempty"`
			}
			t.Run("no data", func(t *testing.T) {
				skipTest(t, "incompatible with std")
				checkEncodeValueWithStandard(t, &st{
					A: testJSONMarshaler2{},
				}, nil)
			})
			t.Run("with data", func(t *testing.T) {
				checkEncodeValueWithStandard(t, &st{
					A: testJSONMarshaler2{
						data: "true",
					},
				}, nil)
			})
		})
		t.Run("pointer field", func(t *testing.T) {
			type st struct {
				A *testJSONMarshaler2 `json:",omitempty"`
			}
			t.Run("nil", func(t *testing.T) {
				skipTest(t, "incompatible with std")
				checkEncodeValueWithStandard(t, &st{}, nil)
			})
			t.Run("no data", func(t *testing.T) {
				skipTest(t, "incompatible with std")
				checkEncodeValueWithStandard(t, &st{
					A: &testJSONMarshaler2{},
				}, nil)
			})
			t.Run("with data", func(t *testing.T) {
				checkEncodeValueWithStandard(t, &st{
					A: &testJSONMarshaler2{
						data: "true",
					},
				}, nil)
			})
		})

	})
	t.Run("dynamic json marshaler", func(t *testing.T) {
		type st struct {
			A json.Marshaler `json:",omitempty"`
		}
		t.Run("nil", func(t *testing.T) {
			checkEncodeValueWithStandard(t, st{}, nil)
		})
		t.Run("nil pointer", func(t *testing.T) {
			checkEncodeValueWithStandard(t, st{
				A: (*testJSONMarshaler2)(nil),
			}, nil)
		})
	})
}
