package commands

import (
	bf "github.com/fvdveen/bf"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("command")

// Commands is the list of all commands
var Commands = []bf.CommandInterface{}

// Register adds the command to Commands
func Register(com bf.CommandInterface) bf.CommandInterface {
	Commands = append(Commands, com)
	return com
}
