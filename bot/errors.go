package bot

import "errors"

var (
	// ErrUnknownVoiceState is used when a users voice state could not be found
	ErrUnknownVoiceState = errors.New("could not find user voice state")
)
