package jzon

type ValueType int

const (
	WhiteSpaceValue ValueType = iota
	InvalidValue
	StringValue
	NumberValue
	ObjectValue
	ArrayValue
	BoolValue
	NullValue
	LastValue
)

var (
	valueTypeNames [LastValue]string
	valueTypeMap   [charNum]ValueType
)

func init() {
	// value type names
	valueTypeNames[WhiteSpaceValue] = "WhiteSpaceValue"
	valueTypeNames[InvalidValue] = "InvalidValue"
	valueTypeNames[StringValue] = "StringValue"
	valueTypeNames[NumberValue] = "NumberValue"
	valueTypeNames[ObjectValue] = "ObjectValue"
	valueTypeNames[ArrayValue] = "ArrayValue"
	valueTypeNames[BoolValue] = "BoolValue"
	valueTypeNames[NullValue] = "NullValue"

	// value type map
	for i := 0; i < charNum; i++ {
		valueTypeMap[i] = InvalidValue
	}
	for _, c := range " \n\t\r" {
		valueTypeMap[c] = WhiteSpaceValue
	}
	valueTypeMap['"'] = StringValue
	for _, c := range "-0123456789" {
		valueTypeMap[c] = NumberValue
	}
	valueTypeMap['{'] = ObjectValue
	valueTypeMap['['] = ArrayValue
	for _, c := range "tf" {
		valueTypeMap[c] = BoolValue
	}
	valueTypeMap['n'] = NullValue
}
