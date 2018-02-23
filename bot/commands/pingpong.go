package commands

import (
	"fmt"

	bf "github.com/FrisovanderVeen/bf"
)

// Ping sends pong to the text channel that triggered it
var Ping = &bf.Command{
	Name:    "ping",
	Trigger: "ping",
	Use:     "Send pong",
	Action: func(ctx bf.Context) {
		if err := ctx.SendMessage("pong"); err != nil {
			ctx.Error(fmt.Errorf("ping: %v", err))
		}
	},
}

// Pong sends ping to the text channel that triggered it
var Pong = &bf.Command{
	Name:    "pong",
	Trigger: "pong",
	Use:     "Send ping",
	Action: func(ctx bf.Context) {
		if err := ctx.SendMessage("ping"); err != nil {
			ctx.Error(fmt.Errorf("pong: %v", err))
		}
	},
}
