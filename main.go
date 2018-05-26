package main

import (
	"github.com/fvdveen/mu2/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		logrus.Fatal(err)
	}
}
