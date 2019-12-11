package jzon

const (
	invalidFloatDigit = -1
	dotInNumber       = -2
	expInNumber       = -3
)

var (
	floatDigits [charNum]int8
)

func init() {
	for i := 0; i < charNum; i++ {
		floatDigits[i] = invalidFloatDigit
	}
	for i := '0'; i <= '9'; i++ {
		floatDigits[i] = int8(i - '0')
	}
	floatDigits['.'] = dotInNumber
	floatDigits['e'] = expInNumber
	floatDigits['E'] = expInNumber
}
