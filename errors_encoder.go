package jzon

import (
	"errors"
)

var ErrNoAttachedWriter = errors.New("no attached writer")

var ErrFloatIsInfinity = errors.New("float is infinity")

var ErrFloatIsNan = errors.New("float is NaN")
