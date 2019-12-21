package jzon

import (
	"errors"
)

var NoWriterAttachedError = errors.New("no writer attached")

var FloatIsInfinity = errors.New("float is infinity")

var FloatIsNan = errors.New("float is NaN")
