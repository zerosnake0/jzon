package jzon

import (
	"testing"
)

type testLoopStruct struct {
	*testLoopStruct
}

type testLoopStruct2 struct {
	*testLoopStruct3
}

type testLoopStruct3 struct {
	*testLoopStruct2
}

type testLoopStruct4 struct {
	A *testLoopStruct4 `json:"a"`
}

type testLoopIFace interface{}

type testLoopStruct5 struct {
	testLoopIFace
}

func TestValDecoder_Native_Struct_Loop(t *testing.T) {
	f := func(t *testing.T, data string, ex error, p1, p2 interface{}) {
		t.Log(">>>>> initValues >>>>>")
		printValue(t, "p1", p1)
		printValue(t, "p2", p2)
		t.Log(">>>>>>>>>>>>>>>>>>>>>>")
		checkDecodeWithStandard(t, DefaultDecoderConfig, data, ex, p1, p2)
		t.Log("<<<<< initValues <<<<<")
		printValue(t, "p1", p1)
		printValue(t, "p2", p2)
		t.Log("<<<<<<<<<<<<<<<<<<<<<<")
	}
	t.Run("self nested", func(t *testing.T) {
		var t1 testLoopStruct
		t1.testLoopStruct = &t1
		var t2 testLoopStruct
		t2.testLoopStruct = &t2
		f(t, "{}", nil, &t1, &t2)
	})
	t.Run("cross nested", func(t *testing.T) {
		var t1 testLoopStruct2
		var t1t testLoopStruct3
		t1.testLoopStruct3 = &t1t
		t1t.testLoopStruct2 = &t1
		var t2 testLoopStruct2
		var t2t testLoopStruct3
		t2.testLoopStruct3 = &t2t
		t2t.testLoopStruct2 = &t2
		f(t, "{}", nil, &t1, &t2)
	})
	t.Run("field nested", func(t *testing.T) {
		var t1 testLoopStruct4
		t1.A = &t1
		var t2 testLoopStruct4
		t2.A = &t2
		f(t, `{"a":{"a":{}}}`, nil, &t1, &t2)
	})
	t.Run("interface nested", func(t *testing.T) {
		var t1 testLoopStruct5
		t1.testLoopIFace = &t1
		var t2 testLoopStruct5
		t2.testLoopIFace = &t2
		f(t, `{}`, nil, &t1, &t2)
	})
}
