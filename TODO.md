# Feature

- [x] Skip object
- [x] Support other tag than `json`
- [x] Decoder option `UseNumber`
- [x] Decoder option `DisallowUnknownFields`
- [x] tag option `quoted`
- [x] json/text marshaler (pointer receiver) for values
- [x] tag option `omitempty`

# Improvement

- [x] Nested skip (by using a stack)
- [x] Decode with stack (proven worse than use recursive directly)

# Benchmark

- [x] Reset
- [x] Skip switch/slice

# Benchmark with other library

- [x] jsoniter
