package pingpong

import (
	bf "github.com/FrisovanderVeen/bf"
	"github.com/FrisovanderVeen/mu2/commands"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("commands/pingpong")

// Ping sends pong to the text channel that triggered it
var _ = commands.Register(bf.NewCommand(
	bf.Name("ping"),
	bf.Trigger("ping"),
	bf.Use("Sends pong to the text channel"),
	bf.Action(func(ctx bf.Context) {
		if err := ctx.SendMessage("pong"); err != nil {
			log.Errorf("%v", err)
		}
	}),
))

// Pong sends ping to the text channel that triggered it
var _ = commands.Register(bf.NewCommand(
	bf.Name("pong"),
	bf.Trigger("pong"),
	bf.Use("Sends ping to the text channel"),
	bf.Action(func(ctx bf.Context) {
		if err := ctx.SendMessage("ping"); err != nil {
			log.Errorf("%v", err)
		}
	}),
))
