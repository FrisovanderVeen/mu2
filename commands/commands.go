package commands

import (
	bf "github.com/fvdveen/bf"
	logging "github.com/op/go-logging"
)

var (
	log     = logging.MustGetLogger("command")
	VERSION string
)

// Commands is the list of all commands
var Commands = []bf.Command{}

// Register adds the command to Commands
func Register(com bf.Command) bf.Command {
	Commands = append(Commands, com)
	return com
}
