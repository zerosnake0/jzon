package jzon

import (
	"testing"
)

type testComplexInner struct {
	testComplexOuter
	A0 int
	A1 int
}

type testComplexInner2 struct {
	B0 int
}

type aliasTestComplexInner2 testComplexInner2

type testComplexInner2Inner struct {
	B0 int
	B1 int
}

type testComplexInner2Outer struct {
	testComplexInner2Inner
}

type testComplexInner3 struct {
	C0 int
}

type renameTestComplexInner3 = testComplexInner3

type testComplexInner3Inner struct {
	C0 int
	C1 int
}

type testComplexInner3Outer struct {
	testComplexInner3Inner
}

type testComplexInner4A struct {
	D0 int // overridden by json tag
	D1 int // the sibling is renamed
	D3 int // the sibling is ignored
	D4 int // duplicated
	D5 int `json:"D6"` // duplicated with json key
	D7 int // duplicated, but inner tagged field will be promoted
}

type testComplexInner4A2 testComplexInner4A

type testComplexInner4A3 testComplexInner4A

type testComplexInner4B struct {
	D0 int `json:"D0"`
	D1 int `json:"D2"`
	D3 int `json:"-"`
	D4 int
	D5 int `json:"D6"`
	D7 int
}

type testComplexInner4CInternal struct {
	D0 int
	D1 int
	D2 int
	D3 int
	D4 int
	D5 int
	D6 int
	D7 int `json:"D7"`
}

type testComplexInner4C struct {
	testComplexInner4CInternal
}

type testComplexInner5A struct {
	E0 int // wins for lesser depth
	E1 int `json:"E1"` // wins for tagged
}

type testComplexInner5BInternal struct {
	E0 int
}

type testComplexInner5B struct {
	testComplexInner5BInternal
	E1 int
}

type testComplexOuter struct {
	// nested
	*testComplexInner
	A0 int

	// alias
	testComplexInner2
	aliasTestComplexInner2
	testComplexInner2Outer

	// duplicated (by rename)
	testComplexInner3
	renameTestComplexInner3
	testComplexInner3Outer

	// duplicate
	testComplexInner4A
	testComplexInner4A2 `json:"e"` // treated as named
	testComplexInner4A3 `json:"-"` // ignored
	testComplexInner4B
	testComplexInner4C

	// Ambigue
	F0 int `json:"AMBIGUE"`
	F1 int `json:"Ambigue"`

	// sort
	testComplexInner5A
	testComplexInner5B
}

func TestValDecoder_Native_Struct_Embedded_Complex(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		checkDecodeWithStandard(t, DefaultDecoder, data, ex, p1, p2)
	}
	t.Run("complex", func(t *testing.T) {
		f(t, `{
			"A0": 10, "A1": 11,
			"B0": 20, "B1": 21,
            "C0": 30, "C1": 31,
			"D0": 40, "D1": 41, "D2": 42,
			"D3": 43, "D4": 44, "D5": 45,
			"D6": 46, "D7": 47,
			"ambigue": 50,
			"E0": 60, "E1": 61 
		}`, nil, &testComplexOuter{
			testComplexInner: &testComplexInner{},
		}, &testComplexOuter{
			testComplexInner: &testComplexInner{},
		})
	})
}
