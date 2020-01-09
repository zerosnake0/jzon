[![Go Report Card](https://goreportcard.com/badge/github.com/zerosnake0/jzon)](https://goreportcard.com/report/github.com/zerosnake0/jzon)
[![Build Status](https://travis-ci.org/zerosnake0/jzon.svg?branch=master)](https://travis-ci.org/zerosnake0/jzon)
[![codecov](https://codecov.io/gh/zerosnake0/jzon/branch/master/graph/badge.svg)](https://codecov.io/gh/zerosnake0/jzon)

# jzon

![](https://github.com/zerosnake0/jzon/workflows/Test/badge.svg)

## Why another jsoniter?

The code I write here is very similar to [github.com/json-iterator/go](https://github.com/json-iterator/go),
so you may ask why reinvent the wheel.

For sure that I benefit a lot from the `jsoniter` library, but i found some inconvenience for me to use it
in some condition, for example:

- the iterator methods ReadString accepts null, there is no method which accepts exactly string.
  I have to do some extra check before calling.
- some behavior is not compatible with the standard library.
- I want a chained streamer

On the other hand, I also want to learn how the `jsoniter` works, so there is this repo.

## What's different from jsoniter?

Here are some of the differences:

- the iterator methods accept the exact type, for example ReadString accepts only string, not null
- the behavior is almost the same as the standard library (when an error returns, the behavior may differ
  from the standard library)
- the error of the iterator is returned instead of being saved inside iterator
- the decoder/encoder interface has additional options, like struct tag options

Some features of `jsoniter` are not implemented, and may be not implemented in the future neither.
I choose only the ones I need to implement.

## Compatibility with standard library

I tried implemented a version which is completely compatible with the standard library:

https://github.com/zerosnake0/jzon/tree/reflect

The benchmark shows that it's much faster than the standard library.
However it is still much slower than the current version,
which cannot be exactly the same as standard library (at least in my POV).

The major incompatibility is about the two following interfaces:
- `json.Marshaler`
- `encoding.TextMarshaler`

The method on pointer receiver may be called with an unaddressable value,
for example:

```go
type field struct {}

func (*field) MarshalJSON() ([]byte, error)

type st struct {
    F field
}

json.Marshal(st{}) // will not call field.MarshalJSON
jzon.Marshal(st{}) // will call field.MarshalJSON
```

So the user should be care when marshaling a value when method on
pointer receiver is involved

You can check the tests for more detailed info about the difference

## How to use

### Standard library like

```go
import "github.com/zerosnake0/jzon"

// Unmarshal
err := jzon.Unmarshal(b, &data)

// Marshal
b, err := jzon.Marshal(&data)
```

### Iterator

```go
iter := jzon.NewIterator()
defer jzon.ReturnIterator(iter)
iter.Reset(b)
jzon.ReadVal(&data)
```

### Streamer

```go
var w io.Writer

streamer := jzon.NewStreamer()
defer jzon.ReturnStreamer(streamer)
streamer.Reset(w)
streamer.Value(&data)
streamer.Flush()
```

### Custom Decoder

see `decoder_test.go`

```go
type testIntDecoder struct{}

func (*testIntDecoder) Decode(ptr unsafe.Pointer, it *Iterator, opts *DecOpts) error {
    ...
}

dec := NewDecoder(&DecoderOption{
    ValDecoders: map[reflect.Type]ValDecoder{
        reflect.TypeOf(int(0)): (*testIntDecoder)(nil),
    },
    CaseSensitive: true,
})

// standard library like
err := dec.Unmarshal(b, &data)

// iterator
iter := dec.NewIterator()
defer dec.ReturnIterator(iter)
```

### Custom Encoder

see `encoder_test.go`

```go
type testIntEncoder struct{}

func (*testIntEncoder) IsEmpty(ptr unsafe.Pointer) bool {
    ...
}

func (*testIntEncoder) Encode(ptr unsafe.Pointer, s *Streamer, opts *EncOpts) {
    ...
}

enc := NewEncoder(&EncoderOption{
    ValEncoders: map[reflect.Type]ValEncoder{
        reflect.TypeOf(int(0)): (*testIntEncoder)(nil),
    },
})

// standard library like
b, err := enc.Marshal(&data)

// streamer
streamer := enc.NewStreamer()
defer enc.ReturnStreamer(streamer)
```