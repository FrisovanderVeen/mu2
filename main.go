package main

import (
	"os"

	"github.com/fvdveen/mu2/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	app := cmd.New()

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
