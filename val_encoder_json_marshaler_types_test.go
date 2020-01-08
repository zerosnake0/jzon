package jzon

import (
	"fmt"
)

// byte
type testMarshalByte byte

func (tb testMarshalByte) MarshalJSON() ([]byte, error) {
	return []byte{'"', '1', byte(tb), '"'}, nil
}

type testMarshalByte2 byte

func (tb *testMarshalByte2) MarshalJSON() ([]byte, error) {
	return []byte{'"', '2', byte(*tb), '"'}, nil
}

type testMarshalByte3 byte

func (tb testMarshalByte3) MarshalText() ([]byte, error) {
	return []byte{'"', '3', byte(tb), '"'}, nil
}

type testMarshalByte4 byte

func (tb *testMarshalByte4) MarshalText() ([]byte, error) {
	return []byte{'"', '4', byte(*tb), '"'}, nil
}

// map
type testMapJsonMarshaler map[int]int

func (m testMapJsonMarshaler) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%d", len(m))
	return []byte(s), nil
}

// struct
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
	if m == nil {
		return []byte(`"is_null"`), nil
		// return []byte(`null`), nil
	}
	return []byte(m.data), m.err
}

/* The following struct definition is not allowed
type testJsonMarshaler3 struct {
}

type pTestJsonMarshaler3 *testJsonMarshaler3

func (pTestJsonMarshaler3) MarshalJSON() ([]byte, error) {
	return nil, nil
}
*/
