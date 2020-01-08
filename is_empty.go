package jzon

import (
	"reflect"
	"strconv"
	"unsafe"
)

type isEmptyFunc func(ptr unsafe.Pointer) bool

var (
	isEmptyFunctions = [numKinds]isEmptyFunc{}
)

func init() {
	isEmptyFunctions[reflect.Bool] = (*boolEncoder)(nil).IsEmpty
	if strconv.IntSize == 32 {
		isEmptyFunctions[reflect.Int] = (*int32Encoder)(nil).IsEmpty
		isEmptyFunctions[reflect.Uint] = (*uint32Encoder)(nil).IsEmpty
	} else {
		isEmptyFunctions[reflect.Int] = (*int64Encoder)(nil).IsEmpty
		isEmptyFunctions[reflect.Uint] = (*uint64Encoder)(nil).IsEmpty
	}
	isEmptyFunctions[reflect.Int8] = (*int8Encoder)(nil).IsEmpty
	isEmptyFunctions[reflect.Int16] = (*int16Encoder)(nil).IsEmpty
	isEmptyFunctions[reflect.Int32] = (*int32Encoder)(nil).IsEmpty
	isEmptyFunctions[reflect.Int64] = (*int64Encoder)(nil).IsEmpty
	isEmptyFunctions[reflect.Uint8] = (*uint8Encoder)(nil).IsEmpty
	isEmptyFunctions[reflect.Uint16] = (*uint16Encoder)(nil).IsEmpty
	isEmptyFunctions[reflect.Uint32] = (*uint32Encoder)(nil).IsEmpty
	isEmptyFunctions[reflect.Uint64] = (*uint64Encoder)(nil).IsEmpty
	if unsafe.Sizeof(uintptr(0)) == 4 {
		isEmptyFunctions[reflect.Uintptr] = (*uint32Encoder)(nil).IsEmpty
	} else {
		isEmptyFunctions[reflect.Uintptr] = (*uint64Encoder)(nil).IsEmpty
	}
	isEmptyFunctions[reflect.Float32] = (*float32Encoder)(nil).IsEmpty
	isEmptyFunctions[reflect.Float64] = (*float64Encoder)(nil).IsEmpty
	isEmptyFunctions[reflect.Array] = (*arrayEncoder)(nil).IsEmpty
	isEmptyFunctions[reflect.Map] = (*directMapEncoder)(nil).IsEmpty
	isEmptyFunctions[reflect.Slice] = (*sliceEncoder)(nil).IsEmpty
	isEmptyFunctions[reflect.Struct] = (*structEncoder)(nil).IsEmpty
	// Interface
	// Ptr
}
