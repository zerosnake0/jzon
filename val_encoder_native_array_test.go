package jzon

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestValEncoder_Array_Error(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		e := errors.New("test")
		arr := [...]json.Marshaler{testJsonMarshaler{
			data: `"test"`,
			err:  e,
		}}
		checkEncodeValueWithStandard(t, arr, e)
	})
}

func TestValEncoder_Array_Empty(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		arr := [...]int{}
		checkEncodeValueWithStandard(t, arr, nil)
	})
	t.Run("pointer", func(t *testing.T) {
		f := func(t *testing.T, ptr *[0]int, err error) {
			checkEncodeValueWithStandard(t, ptr, err)
		}
		t.Run("nil", func(t *testing.T) {
			f(t, nil, nil)
		})
		t.Run("pointer", func(t *testing.T) {
			arr := [...]int{}
			f(t, &arr, nil)
		})
	})
	t.Run("pointer of pointer", func(t *testing.T) {
		f := func(t *testing.T, ptr **[0]int, err error) {
			checkEncodeValueWithStandard(t, ptr, err)
		}
		t.Run("nil", func(t *testing.T) {
			f(t, nil, nil)
		})
		t.Run("pointer of nil", func(t *testing.T) {
			ptr := (*[0]int)(nil)
			f(t, &ptr, nil)
		})
		t.Run("non nil", func(t *testing.T) {
			arr := [...]int{}
			ptr := &arr
			f(t, &ptr, nil)
		})
	})
}

func TestValEncoder_Array_Indirect(t *testing.T) {
	// len != 1
	t.Run("len<>1", func(t *testing.T) {
		t.Run("pointer", func(t *testing.T) {
			arr := [...]int{1, 2, 3}
			checkEncodeValueWithStandard(t, &arr, nil)
		})
		t.Run("non pointer", func(t *testing.T) {
			arr := [...]int{1, 2, 3}
			checkEncodeValueWithStandard(t, arr, nil)
		})
		t.Run("array of pointer", func(t *testing.T) {
			i := 1
			arr := [...]*int{(*int)(nil), &i}
			checkEncodeValueWithStandard(t, arr, nil)
		})
	})
	// element is indirect
	t.Run("len==1", func(t *testing.T) {
		t.Run("one element array", func(t *testing.T) {
			arr := [...]int{1}
			checkEncodeValueWithStandard(t, arr, nil)
		})
	})
}

func TestValEncoder_Array_Direct(t *testing.T) {
	t.Run("nil element", func(t *testing.T) {
		arr := [...]*int{(*int)(nil)}
		checkEncodeValueWithStandard(t, arr, nil)
	})
	t.Run("non nil element", func(t *testing.T) {
		i := 1
		arr := [...]*int{&i}
		checkEncodeValueWithStandard(t, arr, nil)
	})
	t.Run("pointer", func(t *testing.T) {
		checkEncodeValueWithStandard(t, (*[1]*int)(nil), nil)
	})
}

func TestValEncoder_Array_Marshaler(t *testing.T) {
	t.Run("json marshaler", func(t *testing.T) {
		t.Run("arr value", func(t *testing.T) {
			t.Run("value", func(t *testing.T) {
				checkEncodeValueWithStandard(t, [...]testMarshalByte{'t', 'e', 's', 't'}, nil)
			})
			t.Run("pointer", func(t *testing.T) {
				skipTest(t, "array with pointer json marshaler")
				checkEncodeValueWithStandard(t, [...]testMarshalByte2{'t', 'e', 's', 't'}, nil)
			})
		})
		t.Run("arr ptr", func(t *testing.T) {
			t.Run("value", func(t *testing.T) {
				arr := [...]testMarshalByte{'t', 'e', 's', 't'}
				checkEncodeValueWithStandard(t, &arr, nil)
			})
			t.Run("pointer", func(t *testing.T) {
				arr := [...]testMarshalByte2{'t', 'e', 's', 't'}
				checkEncodeValueWithStandard(t, &arr, nil)
			})
		})
	})
	t.Run("text marshaler", func(t *testing.T) {
		t.Run("arr value", func(t *testing.T) {
			t.Run("value", func(t *testing.T) {
				checkEncodeValueWithStandard(t, [...]testMarshalByte3{'t', 'e', 's', 't'}, nil)
			})
			t.Run("pointer", func(t *testing.T) {
				skipTest(t, "array with pointer text marshaler")
				checkEncodeValueWithStandard(t, [...]testMarshalByte4{'t', 'e', 's', 't'}, nil)
			})
		})
		t.Run("arr ptr", func(t *testing.T) {
			t.Run("value", func(t *testing.T) {
				arr := [...]testMarshalByte3{'t', 'e', 's', 't'}
				checkEncodeValueWithStandard(t, &arr, nil)
			})
			t.Run("pointer", func(t *testing.T) {
				arr := [...]testMarshalByte4{'t', 'e', 's', 't'}
				checkEncodeValueWithStandard(t, &arr, nil)
			})
		})
	})
}

func TestValEncoder_Array_OmitEmpty(t *testing.T) {
	t.Run("direct", func(t *testing.T) {
		type st struct {
			A [1]*int `json:",omitempty"`
		}
		checkEncodeValueWithStandard(t, &st{}, nil)
	})
	t.Run("indirect", func(t *testing.T) {
		type st struct {
			A [1]int `json:",omitempty"`
		}
		checkEncodeValueWithStandard(t, &st{}, nil)
	})
}
