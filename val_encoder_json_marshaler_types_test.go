package jzon

import (
	"fmt"
	"strconv"
)

// bool
type testBoolJSONMarshaler bool

func (b testBoolJSONMarshaler) MarshalJSON() ([]byte, error) {
	if b {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}

// array
type testIndirectArrayMarshaler [1]int

func (arr testIndirectArrayMarshaler) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf(`"%d"`, arr[0]*2)
	return []byte(s), nil
}

type testDirectArrayMarshaler [1]*int

func (arr testDirectArrayMarshaler) MarshalJSON() ([]byte, error) {
	if arr[0] == nil {
		return []byte(`"nil"`), nil
	}
	return []byte(strconv.Itoa(*arr[0])), nil
}

// map
type testMapJSONMarshaler map[int]int

func (m testMapJSONMarshaler) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%d", len(m))
	return []byte(s), nil
}

// slice
type testSliceMarshaler []byte

func (s testSliceMarshaler) MarshalJSON() ([]byte, error) {
	str := strconv.Itoa(len(s))
	return []byte(str), nil
}

// struct
type testIndirectStructMarshaler struct {
	A int
}

func (s testIndirectStructMarshaler) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(s.A)), nil
}

type testDirectStructMarshaler struct {
	A *int
}

func (s testDirectStructMarshaler) MarshalJSON() ([]byte, error) {
	if s.A == nil {
		return []byte(`"nil"`), nil
	}
	return []byte(strconv.Itoa(*s.A)), nil
}

type testJSONMarshaler struct {
	data string
	err  error
}

func (m testJSONMarshaler) MarshalJSON() ([]byte, error) {
	return []byte(m.data), m.err
}

type testJSONMarshaler2 struct {
	data string
	err  error
}

func (m *testJSONMarshaler2) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte(`"is_null"`), nil
		// return []byte(`null`), nil
	}
	return []byte(m.data), m.err
}

/* The following struct definition is not allowed
type testJSONMarshaler3 struct {
}

type pTestJSONMarshaler3 *testJSONMarshaler3

func (pTestJSONMarshaler3) MarshalJSON() ([]byte, error) {
	return nil, nil
}
*/
