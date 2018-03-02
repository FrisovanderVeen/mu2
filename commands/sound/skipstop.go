package sound

import (
	"github.com/fvdveen/bf"
	"github.com/fvdveen/mu2/commands"
)

var _ = commands.Register(bf.NewCommand(
	bf.Name("skip"),
	bf.Trigger("skip"),
	bf.Use("Skips the currently playing audio"),
	bf.Action(func(ctx bf.Context) {
		skip <- ctx.Message
	}),
))

var _ = commands.Register(bf.NewCommand(
	bf.Name("stop"),
	bf.Trigger("stop"),
	bf.Use("Clears the queue and stops playing audio"),
	bf.Action(func(ctx bf.Context) {
		stop <- ctx.Message
	}),
))
