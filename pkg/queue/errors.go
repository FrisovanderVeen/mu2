package queue

import "errors"

var (
	// ErrOutOfBounds is used when an given index is out of bounds
	ErrOutOfBounds = errors.New("index out of queue bounds")
)
