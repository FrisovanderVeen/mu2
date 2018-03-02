package main

import (
	"os"

	"github.com/fvdveen/mu2/cmd"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

func main() {
	app := cmd.NewApp()
	if err := app.Run(os.Args); err != nil {
		log.Critical("%v", err)
	}
}
