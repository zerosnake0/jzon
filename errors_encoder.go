package jzon

import (
	"errors"
)

// ErrNoAttachedWriter there is no writer attaching to the streamer
var ErrNoAttachedWriter = errors.New("no attached writer")

// ErrFloatIsInfinity the float to write is infinity
var ErrFloatIsInfinity = errors.New("float is infinity")

// ErrFloatIsNan the float to write is NaN
var ErrFloatIsNan = errors.New("float is NaN")
