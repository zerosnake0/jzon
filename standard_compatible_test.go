package jzon

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"
	"unsafe"

	"github.com/stretchr/testify/require"
)

var (
	compatibleOnError = os.Getenv("COMPATIBLE_ON_ERROR") == "1"
)

func nestedArray1(count int) string {
	return strings.Repeat(" [", count) + " [ ] " +
		strings.Repeat("] ", count)
}

func nestedArray2(count int) string {
	return strings.Repeat(" [ [ ], ", count) + " [ ] " +
		strings.Repeat("] ", count)
}

func nestedArrayWithObject(count int) string {
	return strings.Repeat(" [ { }, ", count) + " [ ] " +
		strings.Repeat("] ", count)
}

func nestedObject(count int) string {
	return strings.Repeat(` { "a" : { }, "b": `, count) + " { } " +
		strings.Repeat("} ", count)
}

func nestedObjectWithArray(count int) string {
	return strings.Repeat(` { "a" : [ ], "b": `, count) + " [ ] " +
		strings.Repeat("} ", count)
}

type printValueKey struct {
	rtype rtype
	ptr   uintptr
}

func (pk printValueKey) String() string {
	return fmt.Sprintf("<%x %x>", pk.rtype, pk.ptr)
}

func printValue(t *testing.T, prefix string, o interface{}) {
	prefix += " "
	if o == nil {
		t.Logf(prefix + "nil")
		return
	}
	visited := map[printValueKey]bool{}
	oV := reflect.ValueOf(o)
	for indent := prefix; ; indent += "  " {
		i := oV.Interface()
		ef := (*eface)(unsafe.Pointer(&i))
		vk := printValueKey{ef.rtype, uintptr(ef.data)}

		k := oV.Kind()
		t.Logf(indent+"%+v %+v %v", oV.Type(), oV, vk)

		if k != reflect.Interface && k != reflect.Ptr {
			break
		}
		if oV.IsNil() {
			break
		}

		if visited[vk] {
			t.Logf(indent + "  visited...")
			break
		}
		visited[vk] = true

		oV = oV.Elem()
	}
}

func checkDecodeWithStandard(t *testing.T, decCfg *DecoderConfig, data string, ex error, exp, got interface{}) {
	b := []byte(data)
	expErr := json.Unmarshal(b, exp)
	gotErr := decCfg.Unmarshal(b, got)
	t.Logf("\nexpErr: %+v\ngotErr: %+v", expErr, gotErr)
	noError := expErr == nil
	if noError {
		printValue(t, "exp", reflect.ValueOf(exp).Elem().Interface())
	}
	require.Equal(t, noError, gotErr == nil,
		"exp %+v\ngot %+v", expErr, gotErr)
	require.Equalf(t, noError, ex == nil, "exp err: %v\ngot err: %v", ex, gotErr)
	if ex != nil {
		checkError(t, ex, gotErr)
		// if reflect.TypeOf(errors.New("")) == reflect.TypeOf(ex) {
		// 	require.Equalf(t, ex, gotErr, "exp err:%v\ngot err:%v", ex, gotErr)
		// } else {
		// 	require.IsTypef(t, ex, gotErr, "exp err:%v\ngot err:%v", ex, gotErr)
		// }
	}
	if !noError && !compatibleOnError {
		return
	}
	if exp == nil {
		require.Equal(t, nil, got)
		return
	}
	expV := reflect.ValueOf(exp)
	gotV := reflect.ValueOf(got)
	if expV.IsNil() {
		require.True(t, gotV.IsNil())
		return
	}
	expI := expV.Elem().Interface()
	gotI := gotV.Elem().Interface()
	printValue(t, "got", gotI)
	require.Equalf(t, expI, gotI, "exp %+v\ngot %+v", expI, gotI)
}

func TestValid(t *testing.T) {
	f := func(t *testing.T, s string) {
		data := localStringToBytes(s)
		require.Equal(t, json.Valid(data), Valid(data))
	}
	t.Run("empty", func(t *testing.T) {
		f(t, "")
	})
	t.Run("empty object", func(t *testing.T) {
		f(t, "{}")
	})
	t.Run("data remained", func(t *testing.T) {
		f(t, "{}1")
	})
}

func TestStandardCompatible_Decoder(t *testing.T) {
	t.Run("consecutive", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			must := require.New(t)

			newReader := func() io.Reader {
				s := ` {} {} `
				return iotest.OneByteReader(strings.NewReader(s))
			}

			check := func(dec decFace, length int, expMore bool, leftOffset, rightOffset int64) {
				buffered := dec.Buffered()
				if buffered == nil {
					return
				}
				b, err := ioutil.ReadAll(buffered)
				must.NoError(err)
				if length > 0 {
					t.Logf("%T: %q", dec, b)
				}
				must.True(length >= len(b))

				offset := dec.InputOffset()
				// t.Logf("%T %d", dec, offset)
				must.True(leftOffset <= offset, "%T %d > %d", dec, leftOffset, offset)
				must.True(rightOffset >= offset, "%T %d < %d", dec, rightOffset, offset)

				more := dec.More()
				must.Equal(expMore, more, "%T", dec)

				offset = dec.InputOffset()
				// t.Logf("%T %d", dec, offset)
				must.True(leftOffset <= offset, "%T %d > %d", dec, leftOffset, offset)
				must.True(rightOffset >= offset, "%T %d < %d", dec, rightOffset, offset)
			}

			f := func(dec decFace) {
				var i, i2, i3 interface{}

				check(dec, 0, true, 0, 1)

				err := dec.Decode(&i)
				must.NoError(err, "%T", dec)
				check(dec, 0, true, 3, 4)

				err2 := dec.Decode(&i2)
				must.NoError(err2, "%T", dec)
				check(dec, 0, false, 6, 7)

				err3 := dec.Decode(&i3)
				must.Equal(io.EOF, err3, "%T", dec)
				check(dec, 1, false, 6, 7)
			}
			f(json.NewDecoder(newReader()))
			f(NewDecoder(newReader()))
		})
		t.Run("failure at start", func(t *testing.T) {
			must := require.New(t)

			newReader := func() io.Reader {
				s := ` } {} `
				return iotest.OneByteReader(strings.NewReader(s))
			}
			f := func(dec decFace) {
				var i, i2 interface{}

				err := dec.Decode(&i)
				t.Logf("%T %v", dec, err)
				must.Error(err)

				err2 := dec.Decode(&i2)
				t.Logf("%T %v", dec, err2)
				must.Equal(err2, err)
			}
			f(json.NewDecoder(newReader()))
			f(NewDecoder(newReader()))
		})
		t.Run("failure at middle", func(t *testing.T) {
			must := require.New(t)

			newReader := func() io.Reader {
				s := ` {} } `
				return iotest.OneByteReader(strings.NewReader(s))
			}
			f := func(dec decFace) {
				var i, i2 interface{}

				err := dec.Decode(&i)
				must.NoError(err)

				err2 := dec.Decode(&i2)
				t.Logf("%T %v", dec, err2)
				must.Error(err2)
			}
			f(json.NewDecoder(newReader()))
			f(NewDecoder(newReader()))
		})
	})
}
