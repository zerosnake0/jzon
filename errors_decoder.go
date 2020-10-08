package jzon

import (
	"errors"
	"fmt"
)

// DecodeError describes the encountered error and its location.
type DecodeError struct {
	reason   error
	location string
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("%s (near %s)", e.reason.Error(), e.location)
}

// ErrDataRemained there is still data remained in the buffer
// normally returned by Unmarshal methods
var ErrDataRemained = errors.New("expecting EOF, but there is still data")

// ErrPointerReceiver if the decode target is not a pointer
var ErrPointerReceiver = errors.New("the receiver is not a pointer")

// ErrNilPointerReceiver if the decode target is nil
var ErrNilPointerReceiver = errors.New("the receiver is nil")

// ErrEmptyIFace if the decode target is an empty interface (with method)
var ErrEmptyIFace = errors.New("cannot unmarshal on empty iface")

// ErrNilEmbeddedPointer if the decode target has an unexported nil field
var ErrNilEmbeddedPointer = errors.New("cannot unmarshal on nil pointer (unexported embedded)")

// ErrEfaceLooping if we encountered a loop
// for example:
//   type iface interface{}
//   var o1 iface
//   o1 = &o1
//   err := DefaultDecoderConfig.Unmarshal([]byte(`1`), o1)
var ErrEfaceLooping = errors.New("eface looping detected")

// InvalidStringCharError there is an invalid character when reading string
type InvalidStringCharError struct {
	c byte
}

func (e InvalidStringCharError) Error() string {
	return fmt.Sprintf("invalid character %x found", e.c)
}

// InvalidEscapeCharError there is an invalid escape character (when reading string)
type InvalidEscapeCharError struct {
	c byte
}

func (e InvalidEscapeCharError) Error() string {
	return fmt.Sprintf("invalid escape character \\%x found", e.c)
}

// InvalidUnicodeCharError there is an invalid unicode character (when reading string)
type InvalidUnicodeCharError struct {
	c byte
}

func (e InvalidUnicodeCharError) Error() string {
	return fmt.Sprintf("invalid unicode character %x found", e.c)
}

// UnexpectedByteError there is an unexpected character
type UnexpectedByteError struct {
	got  byte
	exp  byte
	exp2 byte
}

func (e UnexpectedByteError) Error() string {
	if e.exp == 0 {
		return fmt.Sprintf("unexpected character %q", e.got)
	}
	if e.exp2 == 0 {
		return fmt.Sprintf("expecting %q but got %q", e.exp, e.got)
	}
	return fmt.Sprintf("expecting %q or %q but got %q", e.exp, e.exp2, e.got)
}

// IntOverflowError the integer overflows.
type IntOverflowError struct {
	typ   string
	value string
}

func (e IntOverflowError) Error() string {
	return fmt.Sprintf("overflow %s: %s", e.typ, e.value)
}

// InvalidDigitError there is an invalid digit when reading number
type InvalidDigitError struct {
	c byte
}

func (e InvalidDigitError) Error() string {
	return fmt.Sprintf("invalid digit character: %q", e.c)
}

// InvalidFloatError there is an invalid digit when reading float
type InvalidFloatError struct {
	c byte
}

func (e InvalidFloatError) Error() string {
	return fmt.Sprintf("invalid float character: %q", e.c)
}

// TypeNotSupportedError the decode target is not supported
type TypeNotSupportedError string

func (e TypeNotSupportedError) Error() string {
	return fmt.Sprintf("%q is not supported", string(e))
}

// UnknownFieldError there is an unknown field when decoding
type UnknownFieldError string

func (e UnknownFieldError) Error() string {
	return fmt.Sprintf("unknown field %q", string(e))
}

// BadQuotedStringError the value of field is not correctly quoted
type BadQuotedStringError string

func (e BadQuotedStringError) Error() string {
	return fmt.Sprintf("bad quoted string %q", string(e))
}
