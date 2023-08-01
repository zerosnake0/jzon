package jzon

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValEncoder_Iface_Error(t *testing.T) {
	t.Run("chain error", func(t *testing.T) {
		t.Run("", func(t *testing.T) {
			testStreamerChainError(t, func(s *Streamer) {
				(*ifaceValEncoder)(nil).Encode(nil, s, nil)
			})
		})
	})
}

type testBarface interface {
	Bar() (string, error)
}

type testBarfaceImpl struct {
	s string
	e error
}

func (i testBarfaceImpl) Bar() (string, error) {
	return i.s, i.e
}

type testBarfaceImpl2 struct {
	s string
	e error
}

func (i *testBarfaceImpl2) Bar() (string, error) {
	return i.s, i.e
}

type testBarfaceEncoder struct{}

func (t testBarfaceEncoder) Encode(o interface{}, s *Streamer, opts *EncOpts) {
	if s.Error != nil {
		return
	}
	if o == nil {
		s.Null()
		return
	}
	bface, ok := o.(testBarface)
	if !ok {
		panic("should not reach")
	} else {
		b, err := bface.Bar()
		if err != nil {
			s.Error = err
			return
		}
		s.String(b)
	}
}

func TestValEncoder_Iface_NonPointerReceiver(t *testing.T) {
	cfg := NewEncoderConfig(&EncoderOption{
		IfaceEncoders: []IfaceValEncoderConfig{{
			Type:    reflect.TypeOf((*testBarface)(nil)).Elem(),
			Encoder: testBarfaceEncoder{},
		}},
	})
	t.Run("non pointer", func(t *testing.T) {
		t.Run("no error", func(t *testing.T) {
			m := require.New(t)
			s := strconv.Itoa(rand.Int())
			b, err := cfg.Marshal(testBarfaceImpl{
				s: s,
			})
			m.NoError(err)
			m.Equal(strconv.Quote(s), string(b))
		})
	})
	t.Run("pointer", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			m := require.New(t)
			b, err := cfg.Marshal((*testBarfaceImpl)(nil))
			m.NoError(err)
			m.Equal("null", string(b))
		})
		t.Run("no error", func(t *testing.T) {
			m := require.New(t)
			s := strconv.Itoa(rand.Int())
			b, err := cfg.Marshal(&testBarfaceImpl{
				s: s,
			})
			m.NoError(err)
			m.Equal(strconv.Quote(s), string(b))
		})
	})
	t.Run("pointer of pointer", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			m := require.New(t)
			b, err := cfg.Marshal((**testBarfaceImpl)(nil))
			m.NoError(err)
			m.Equal("null", string(b))
		})
		t.Run("pointer of nil", func(t *testing.T) {
			m := require.New(t)
			ptr := (*testBarfaceImpl)(nil)
			b, err := cfg.Marshal(&ptr)
			m.NoError(err)
			m.Equal("null", string(b))
		})
		t.Run("no error", func(t *testing.T) {
			m := require.New(t)
			s := strconv.Itoa(rand.Int())
			ptr := &testBarfaceImpl{
				s: s,
			}
			b, err := cfg.Marshal(&ptr)
			m.NoError(err)
			m.Equal(strconv.Quote(s), string(b))
		})
		t.Run("error", func(t *testing.T) {
			// TODO
		})
	})
}

func TestValEncoder_Iface_PointerReceiver(t *testing.T) {
	cfg := NewEncoderConfig(&EncoderOption{
		IfaceEncoders: []IfaceValEncoderConfig{{
			Type:    reflect.TypeOf((*testBarface)(nil)).Elem(),
			Encoder: testBarfaceEncoder{},
		}},
	})
	t.Run("non pointer", func(t *testing.T) {
		t.Run("no error", func(t *testing.T) {
			m := require.New(t)
			b, err := cfg.Marshal(testBarfaceImpl2{})
			m.NoError(err)
			m.Equal("{}", string(b))
		})
	})
	t.Run("pointer", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			m := require.New(t)
			b, err := cfg.Marshal((*testBarfaceImpl2)(nil))
			m.NoError(err)
			m.Equal("null", string(b))
		})
		t.Run("no error", func(t *testing.T) {
			m := require.New(t)
			s := strconv.Itoa(rand.Int())
			b, err := cfg.Marshal(&testBarfaceImpl2{
				s: s,
			})
			m.NoError(err)
			m.Equal(strconv.Quote(s), string(b))
		})
	})
	t.Run("pointer of pointer", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			m := require.New(t)
			b, err := cfg.Marshal((**testBarfaceImpl2)(nil))
			m.NoError(err)
			m.Equal("null", string(b))
		})
		t.Run("pointer of nil", func(t *testing.T) {
			m := require.New(t)
			ptr := (*testBarfaceImpl2)(nil)
			b, err := cfg.Marshal(&ptr)
			m.NoError(err)
			m.Equal("null", string(b))
		})
		t.Run("no error", func(t *testing.T) {
			m := require.New(t)
			s := strconv.Itoa(rand.Int())
			ptr := &testBarfaceImpl2{
				s: s,
			}
			b, err := cfg.Marshal(&ptr)
			m.NoError(err)
			m.Equal(strconv.Quote(s), string(b))
		})
		t.Run("error", func(t *testing.T) {
			// TODO
		})
	})
	t.Run("struct member", func(t *testing.T) {
		t.Run("value", func(t *testing.T) {
			type st struct {
				A testBarfaceImpl2
			}
			t.Run("value", func(t *testing.T) {
				m := require.New(t)
				s := strconv.Itoa(rand.Int())
				b, err := cfg.Marshal(st{
					A: testBarfaceImpl2{
						s: s,
					},
				})
				m.NoError(err)
				m.Equal(`{"A":"`+s+`"}`, string(b))
			})
			t.Run("ptr", func(t *testing.T) {
				m := require.New(t)
				s := strconv.Itoa(rand.Int())
				b, err := cfg.Marshal(&st{
					A: testBarfaceImpl2{
						s: s,
					},
				})
				m.NoError(err)
				m.Equal(`{"A":"`+s+`"}`, string(b))
			})
		})
		t.Run("pointer", func(t *testing.T) {
			type st struct {
				A *testBarfaceImpl2
			}
			t.Run("nil", func(t *testing.T) {
				m := require.New(t)
				b, err := cfg.Marshal(&st{})
				m.NoError(err)
				m.Equal(`{"A":null}`, string(b))
			})
		})
	})
}

func TestValEncoder_Iface_Dynamic(t *testing.T) {
	cfg := NewEncoderConfig(&EncoderOption{
		IfaceEncoders: []IfaceValEncoderConfig{{
			Type:    reflect.TypeOf((*testBarface)(nil)).Elem(),
			Encoder: testBarfaceEncoder{},
		}},
	})
	t.Run("iface nil", func(t *testing.T) {
		m := require.New(t)
		var i testBarface
		b, err := cfg.Marshal(&i)
		m.NoError(err)
		m.Equal("null", string(b))
	})
	t.Run("iface", func(t *testing.T) {
		m := require.New(t)
		s := strconv.Itoa(rand.Int())
		var i testBarface = testBarfaceImpl{
			s: s,
		}
		b, err := cfg.Marshal(&i)
		m.NoError(err)
		m.Equal(strconv.Quote(s), string(b))
	})
	t.Run("iface 2", func(t *testing.T) {
		m := require.New(t)
		s := strconv.Itoa(rand.Int())
		var i testBarface = &testBarfaceImpl2{
			s: s,
		}
		b, err := cfg.Marshal(&i)
		m.NoError(err)
		m.Equal(strconv.Quote(s), string(b))
	})
}

type testMapBarface map[int]int

func (m testMapBarface) Bar() (string, error) {
	s := fmt.Sprintf("%d", len(m))
	return s, nil
}

func TestValEncoder_Iface_Direct(t *testing.T) {
	cfg := NewEncoderConfig(&EncoderOption{
		IfaceEncoders: []IfaceValEncoderConfig{{
			Type:    reflect.TypeOf((*testBarface)(nil)).Elem(),
			Encoder: testBarfaceEncoder{},
		}},
	})
	t.Run("value", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			m := require.New(t)
			b, err := cfg.Marshal(testMapBarface(nil))
			m.NoError(err)
			m.Equal(`"0"`, string(b))
		})
		t.Run("non nil", func(t *testing.T) {
			m := require.New(t)
			b, err := cfg.Marshal(testMapBarface{
				1: 2,
			})
			m.NoError(err)
			m.Equal(`"1"`, string(b))
		})
	})
	t.Run("pointer", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			m := require.New(t)
			b, err := cfg.Marshal((*testMapBarface)(nil))
			m.NoError(err)
			m.Equal(`null`, string(b))
		})
		t.Run("non nil", func(t *testing.T) {
			m := require.New(t)
			var i testMapBarface
			b, err := cfg.Marshal(&i)
			m.NoError(err)
			m.Equal(`"0"`, string(b))
		})
		t.Run("non nil 2", func(t *testing.T) {
			m := require.New(t)
			b, err := cfg.Marshal(&testMapBarface{
				1: 2,
			})
			m.NoError(err)
			m.Equal(`"1"`, string(b))
		})
	})
	t.Run("struct member", func(t *testing.T) {
		type st struct {
			A testMapBarface
		}
		m := require.New(t)
		b, err := cfg.Marshal(&st{})
		m.NoError(err)
		m.Equal(`{"A":"0"}`, string(b))
	})
	t.Run("value of iface", func(t *testing.T) {
		t.Run("value", func(t *testing.T) {
			var bar testBarface = testMapBarface{
				1: 2,
			}
			m := require.New(t)
			b, err := cfg.Marshal(bar)
			m.NoError(err)
			m.Equal(`"1"`, string(b))
		})
		t.Run("pointer", func(t *testing.T) {
			var bar testBarface = &testMapBarface{
				1: 2,
			}
			m := require.New(t)
			b, err := cfg.Marshal(bar)
			m.NoError(err)
			m.Equal(`"1"`, string(b))
		})
	})
	t.Run("pointer of iface", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			var bar testBarface
			m := require.New(t)
			b, err := cfg.Marshal(&bar)
			m.NoError(err)
			m.Equal(`null`, string(b))
		})
		t.Run("value", func(t *testing.T) {
			var bar testBarface = testMapBarface{
				1: 2,
			}
			m := require.New(t)
			b, err := cfg.Marshal(&bar)
			m.NoError(err)
			m.Equal(`"1"`, string(b))
		})
		t.Run("pointer", func(t *testing.T) {
			var bar testBarface = &testMapBarface{
				1: 2,
			}
			m := require.New(t)
			b, err := cfg.Marshal(&bar)
			m.NoError(err)
			m.Equal(`"1"`, string(b))
		})
	})
}

func TestValEncoder_Iface_OmitEmpty(t *testing.T) {
	cfg := NewEncoderConfig(&EncoderOption{
		IfaceEncoders: []IfaceValEncoderConfig{{
			Type:    reflect.TypeOf((*testBarface)(nil)).Elem(),
			Encoder: testBarfaceEncoder{},
		}},
	})
	t.Run("", func(t *testing.T) {
		type st struct {
			A testBarfaceImpl `json:",omitempty"`
		}
		m := require.New(t)
		s := strconv.Itoa(rand.Int())
		b, err := cfg.Marshal(st{
			A: testBarfaceImpl{
				s: s,
			},
		})
		m.NoError(err)
		m.Equal(`{"A":"`+s+`"}`, string(b))
	})
}
