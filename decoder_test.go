package jzon

import (
	"reflect"
	"runtime/debug"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

type testIntDecoder struct{}

func (*testIntDecoder) Decode(ptr unsafe.Pointer, it *Iterator, opts *DecOpts) error {
	vt, err := it.NextValueType()
	if err != nil {
		return err
	}
	switch vt {
	case NullValue:
		if err := it.ReadNull(); err != nil {
			return err
		}
		*(*int)(ptr) = 1
		return nil
	default:
		i, err := it.ReadInt()
		if err != nil {
			return err
		}
		*(*int)(ptr) = i * 2
		return nil
	}
}

type testPtrDecoder struct{}

func (*testPtrDecoder) Decode(ptr unsafe.Pointer, it *Iterator, opts *DecOpts) error {
	vt, err := it.NextValueType()
	if err != nil {
		return err
	}
	switch vt {
	case NullValue:
		if err := it.ReadNull(); err != nil {
			return err
		}
		v := -1
		*(**int)(ptr) = &v
		return nil
	default:
		i, err := it.ReadInt()
		if err != nil {
			return err
		}
		i = i * 2
		*(**int)(ptr) = &i
		return nil
	}
}

func TestDecoder_CustomDecoder(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		dec := NewDecoder(&DecoderOption{
			ValDecoders: map[reflect.Type]ValDecoder{
				reflect.TypeOf(int(0)): (*testIntDecoder)(nil),
			},
			CaseSensitive: true,
		})
		t.Run("value", func(t *testing.T) {
			t.Run("null", func(t *testing.T) {
				var i int
				err := dec.UnmarshalFromString("null", &i)
				require.NoError(t, err)
				require.Equal(t, 1, i)
			})
			t.Run("not null", func(t *testing.T) {
				var i int
				err := dec.UnmarshalFromString("1", &i)
				require.NoError(t, err)
				require.Equal(t, 2, i)
			})
		})
		t.Run("pointer", func(t *testing.T) {
			t.Run("null on not null", func(t *testing.T) {
				v := 1
				i := &v
				err := dec.UnmarshalFromString("null", &i)
				require.NoError(t, err)
				require.Nil(t, i)
			})
			t.Run("not null", func(t *testing.T) {
				v := -1
				i := &v
				err := dec.UnmarshalFromString("2", &i)
				require.NoError(t, err)
				require.NotNil(t, i)
				require.Equal(t, 4, v)
			})
		})
		t.Run("field", func(t *testing.T) {
			t.Run("value", func(t *testing.T) {
				st := struct {
					A int `json:"a"`
				}{}
				t.Run("not present", func(t *testing.T) {
					v := 1
					st.A = v
					err := dec.UnmarshalFromString(`{}`, &st)
					require.NoError(t, err)
					require.Equal(t, v, st.A)
				})
				t.Run("null", func(t *testing.T) {
					st.A = -1
					err := dec.UnmarshalFromString(`{"a":null}`, &st)
					require.NoError(t, err)
					require.Equal(t, 1, st.A)
				})
				t.Run("not null", func(t *testing.T) {
					st.A = -1
					err := dec.UnmarshalFromString(`{"a":1}`, &st)
					require.NoError(t, err)
					require.Equal(t, 2, st.A)
				})
			})
			t.Run("pointer", func(t *testing.T) {
				st := struct {
					A *int `json:"a"`
				}{}
				t.Run("not present", func(t *testing.T) {
					v := 123
					st.A = &v
					err := dec.UnmarshalFromString(`{}`, &st)
					require.NoError(t, err)
					require.Equal(t, &v, st.A)
				})
				t.Run("null on not null", func(t *testing.T) {
					v := 123
					st.A = &v
					err := dec.UnmarshalFromString(`{"a":null}`, &st)
					require.NoError(t, err)
					require.Nil(t, st.A)
				})
				t.Run("not null", func(t *testing.T) {
					v := 123
					st.A = &v
					err := dec.UnmarshalFromString(`{"a":1}`, &st)
					require.NoError(t, err)
					require.NotNil(t, st.A)
					require.Equal(t, 2, *st.A)
				})
			})
		})
	})
	t.Run("ptr", func(t *testing.T) {
		dec := NewDecoder(&DecoderOption{
			ValDecoders: map[reflect.Type]ValDecoder{
				reflect.TypeOf((*int)(nil)): (*testPtrDecoder)(nil),
			},
			CaseSensitive: true,
		})
		t.Run("value", func(t *testing.T) {
			t.Run("null", func(t *testing.T) {
				var i *int
				err := dec.UnmarshalFromString("null", &i)
				require.NoError(t, err)
				require.NotNil(t, i)
				require.Equal(t, -1, *i)
			})
			t.Run("not null", func(t *testing.T) {
				var i *int
				err := dec.UnmarshalFromString("1", &i)
				require.NoError(t, err)
				require.NotNil(t, i)
				require.Equal(t, 2, *i)
			})
		})
		t.Run("struct", func(t *testing.T) {
			t.Run("value", func(t *testing.T) {
				st := struct {
					A *int `json:"a"`
				}{}
				t.Run("null", func(t *testing.T) {
					i := 123
					st.A = &i
					err := dec.UnmarshalFromString("null", &st)
					require.NoError(t, err)
					require.Equal(t, &i, st.A)
				})
				t.Run("not present", func(t *testing.T) {
					i := 123
					st.A = &i
					err := dec.UnmarshalFromString(`{}`, &st)
					require.NoError(t, err)
					require.Equal(t, &i, st.A)
				})
				t.Run("present as null", func(t *testing.T) {
					i := 123
					st.A = &i
					err := dec.UnmarshalFromString(`{"a":null}`, &st)
					require.NoError(t, err)
					require.NotNil(t, st.A)
					require.Equal(t, -1, *st.A)
				})
				t.Run("present as not null", func(t *testing.T) {
					i := 123
					st.A = &i
					err := dec.UnmarshalFromString(`{"a":1}`, &st)
					require.NoError(t, err)
					require.NotNil(t, st.A)
					require.Equal(t, 2, *st.A)
				})
			})
			t.Run("pointer", func(t *testing.T) {
				st := struct {
					A **int `json:"a"`
				}{}
				var v *int
				st.A = &v
				require.NotNil(t, st.A)
				err := dec.UnmarshalFromString(`{"a":null}`, &st)
				require.NoError(t, err)
				require.Nil(t, st.A)
			})
		})
	})
	debug.FreeOSMemory()
}
