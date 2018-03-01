package bf

import "errors"

var (
	// ErrVSNotFound  is used when a voice state can't be found
	ErrVSNotFound = errors.New("Could not find user's voice state")
	// ErrDoubleCommand is used when 2 or more commands with the same trigger are added
	ErrDoubleCommand = errors.New("2 or more commands with the same trigger detected, only one of them will never be called")
)
