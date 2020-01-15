package jzon

import (
	"errors"
	"fmt"
)

type DecodeError struct {
	reason   error
	location string
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("%s (near %s)", e.reason.Error(), e.location)
}

// DataRemainedError
var DataRemainedError = errors.New("expecting EOF, but there is still data")

// PointerReceiverError
var PointerReceiverError = errors.New("the receiver is not a pointer")

// NilPointerReceiverError
var NilPointerReceiverError = errors.New("the receiver is nil")

// IFaceError
var IFaceError = errors.New("cannot unmarshal on empty iface")

// NilEmbeddedError
var NilEmbeddedPointerError = errors.New("cannot unmarshal on nil pointer (unexported embedded)")

// EfaceLoopingError
var EfaceLoopingError = errors.New("eface looping detected")

// InvalidStringCharError
type InvalidStringCharError struct {
	c byte
}

func (e InvalidStringCharError) Error() string {
	return fmt.Sprintf("invalid character %x found", e.c)
}

// InvalidEscapeCharError
type InvalidEscapeCharError struct {
	c byte
}

func (e InvalidEscapeCharError) Error() string {
	return fmt.Sprintf("invalid escape character \\%x found", e.c)
}

// InvalidUnicodeCharError
type InvalidUnicodeCharError struct {
	c byte
}

func (e InvalidUnicodeCharError) Error() string {
	return fmt.Sprintf("invalid unicode character %x found", e.c)
}

// UnexpectedByteError
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

// IntOverflow
type IntOverflowError struct {
	typ   string
	value string
}

func (e IntOverflowError) Error() string {
	return fmt.Sprintf("overflow %s: %s", e.typ, e.value)
}

// InvalidDigit
type InvalidDigitError struct {
	c byte
}

func (e InvalidDigitError) Error() string {
	return fmt.Sprintf("invalid digit character: %q", e.c)
}

// InvalidFloat
type InvalidFloatError struct {
	c byte
}

func (e InvalidFloatError) Error() string {
	return fmt.Sprintf("invalid float character: %q", e.c)
}

// TypeNotSupported
type TypeNotSupportedError string

func (e TypeNotSupportedError) Error() string {
	return fmt.Sprintf("%q is not supported", string(e))
}

// UnknownField
type UnknownFieldError string

func (e UnknownFieldError) Error() string {
	return fmt.Sprintf("unknown field %q", string(e))
}

// BadQuotedString
type BadQuotedStringError string

func (e BadQuotedStringError) Error() string {
	return fmt.Sprintf("bad quoted string %q", string(e))
}
