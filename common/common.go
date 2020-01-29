package common

import (
	"io"

	"github.com/fvdveen/mu2/config"
)

var (
	conf config.Config
)

// ReadConfig reads the config from r
func ReadConfig(r io.Reader) error {
	c, err := config.Unmarshal(r)
	if err != nil {
		return err
	}

	conf = *c

	return nil
}

// GetConfig returns the config
func GetConfig() config.Config {
	return conf
}

// GetCommandConfig unmarshals the config for command com into i
func GetCommandConfig(com string, i interface{}) error {
	return conf.UnmarshalCommandConfig(com, i)
}

// GetPrefix returns the prefix of the bot
func GetPrefix() string {
	return conf.Bot.CommandPrefix
}

type commandError struct {
	Command string
	Err     error
}

func (err *commandError) Error() string {
	return "[" + err.Command + "] " + err.Err.Error()
}

func (err *commandError) Unwrap() error {
	return err.Err
}

func CommandError(com string, err error) error {
	if err == nil {
		return nil
	}

	return &commandError{
		Command: com,
		Err:     err,
	}
}

func IsCommandError(err error) (string, bool) {
	com, ok := err.(*commandError)
	if !ok {
		return "", false
	}

	return com.Command, true
}
