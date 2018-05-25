package main

import (
	"github.com/fvdveen/mu2/cmd"
	logging "github.com/op/go-logging"
)

var (
	log     = logging.MustGetLogger("main")
	VERSION = "0.1.1"
)

func main() {
	if err := cmd.Execute(VERSION); err != nil {
		log.Critical(err.Error())
	}
}
