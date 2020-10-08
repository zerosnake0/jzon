package jzon

// ValueType is the type of the next json token
type ValueType int

const (
	// WhiteSpaceValue the next token is whitespace
	WhiteSpaceValue ValueType = iota
	// InvalidValue an error occurred
	InvalidValue
	// StringValue the next token is string
	StringValue
	// NumberValue the next token is number
	NumberValue
	// ObjectValue the next token is object
	ObjectValue
	// ArrayValue the next token is array
	ArrayValue
	// BoolValue the next token is a boolean value
	BoolValue
	// NullValue the next token is null
	NullValue
	// LastValue is a counter, should not be used
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
