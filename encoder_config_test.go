package jzon

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

type testIntEncoder struct{}

func (*testIntEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *(*int)(ptr) == 1
}

func (*testIntEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
	s.Int(*(*int)(ptr) + 1)
}

func TestEncoderConfig_CustomConfig(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		encCfg := NewEncoderConfig(&EncoderOption{
			ValEncoders: map[reflect.Type]ValEncoder{
				reflect.TypeOf(int(0)): (*testIntEncoder)(nil),
			},
		})
		t.Run("value", func(t *testing.T) {
			b, err := encCfg.Marshal(1)
			require.NoError(t, err)
			require.Equal(t, "2", string(b))
		})
		t.Run("pointer", func(t *testing.T) {
			i := 1
			b, err := encCfg.Marshal(&i)
			require.NoError(t, err)
			require.Equal(t, "2", string(b))
		})
		t.Run("struct", func(t *testing.T) {
			t.Run("value", func(t *testing.T) {
				type st struct {
					I int `json:",omitempty"`
				}
				t.Run("zero", func(t *testing.T) {
					b, err := encCfg.Marshal(st{})
					require.NoError(t, err)
					require.Equal(t, `{"I":1}`, string(b))
				})
				t.Run("empty", func(t *testing.T) {
					b, err := encCfg.Marshal(st{I: 1})
					require.NoError(t, err)
					require.Equal(t, `{}`, string(b))
				})
			})
			t.Run("pointer", func(t *testing.T) {
				type st struct {
					I *int `json:",omitempty"`
				}
				t.Run("nil", func(t *testing.T) {
					b, err := encCfg.Marshal(st{})
					require.NoError(t, err)
					require.Equal(t, `{}`, string(b))
				})
				t.Run("zero", func(t *testing.T) {
					i := 0
					b, err := encCfg.Marshal(st{I: &i})
					require.NoError(t, err)
					require.Equal(t, `{"I":1}`, string(b))
				})
				t.Run("empty", func(t *testing.T) {
					i := 1
					b, err := encCfg.Marshal(st{I: &i})
					require.NoError(t, err)
					// pointer is not nil so it's not considered as empty
					require.Equal(t, `{"I":2}`, string(b))
				})
			})
		})
	})
}
