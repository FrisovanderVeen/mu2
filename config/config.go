package config

import (
	"encoding/json"
	"io"
	"strconv"
)

type Config struct {
	Bot      Bot                        `json:"bot"`
	Logger   Logger                     `json:"log"`
	Commands map[string]json.RawMessage `json:"commands"`
	Service  Service                    `json:"service"`
}

type Bot struct {
	Token         string `json:"token"`
	CommandPrefix string `json:"command-prefix"`
}

type Logger struct {
	Level string `json:"level"`
}

type Service struct {
	Youtube Youtube `json:"youtube"`
}

type Youtube struct {
	APIKey string `json:"api-key"`
}

type ErrorType uint8

const (
	ErrTypeNoSuchCommand ErrorType = iota
	ErrUnmarshalCommand
	ErrDecode
)

type Error struct {
	Type  ErrorType
	Key   string
	Cause error
}

func (e *Error) Error() string {
	var msg string

	switch e.Type {
	case ErrTypeNoSuchCommand:
		msg = "no such command: " + e.Key
	case ErrUnmarshalCommand:
		if e.Cause != nil {
			msg = "could not unmarshal command: " + e.Key + ": " + e.Cause.Error()
		} else {
			msg = "could not unmarshal command: " + e.Key
		}
	case ErrDecode:
		if e.Cause != nil {
			msg = "could not decode: " + e.Cause.Error()
		} else {
			msg = "could not decode"
		}
	default:
		msg = "(unknown error type: " + strconv.Itoa(int(e.Type)) + ")"
	}

	return msg
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func Unmarshal(r io.Reader) (*Config, error) {
	c := &Config{
		Commands: map[string]json.RawMessage{},
	}

	err := json.NewDecoder(r).Decode(c)
	if err != nil {
		return nil, &Error{
			Type:  ErrDecode,
			Cause: err,
		}
	}

	return c, nil
}

func (c *Config) UnmarshalCommandConfig(com string, i interface{}) error {
	rm, ok := c.Commands[com]
	if !ok {
		return &Error{
			Type: ErrTypeNoSuchCommand,
			Key:  com,
		}
	}

	err := json.Unmarshal(rm, i)
	if err != nil {
		return &Error{
			Type:  ErrUnmarshalCommand,
			Key:   com,
			Cause: err,
		}
	}

	return nil
}
