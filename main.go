package main

import (
	"os"

	"github.com/FrisovanderVeen/mu2/cmd"
	"github.com/Sirupsen/logrus"
)

func main() {
	app := cmd.NewApp()
	if err := app.Run(os.Args); err != nil {
		logrus.Errorln(err)
	}
}
