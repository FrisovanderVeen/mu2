package commands

import (
	"fmt"

	"github.com/fvdveen/mu2/bot"
)

var (
	commands = make(map[string]bot.Command)
)

// Register registers a command
func Register(c bot.Command) bot.Command {
	commands[c.Name()] = c
	return c
}

// Get returns a command if found
func Get(name string) (bot.Command, error) {
	c, ok := commands[name]
	if !ok {
		return nil, fmt.Errorf("command not found: %s", name)
	}

	return c, nil
}

// All returns all commands
func All() []bot.Command {
	var cs []bot.Command
	for _, c := range commands {
		cs = append(cs, c)
	}

	if cs == nil {
		cs = []bot.Command{}
	}
	return cs
}
