package pingpong

import (
	"github.com/fvdveen/mu2/bot"
	"github.com/fvdveen/mu2/commands"
	"github.com/sirupsen/logrus"
)

var _ = commands.Register(bot.NewCommand("ping", "sends pong", func(ctx bot.Context, args []string) {
	if err := ctx.Send("pong"); err != nil {
		logrus.WithField("command", "ping").Errorf("Send message: %v", err)
	}
}))

var _ = commands.Register(bot.NewCommand("pong", "sends ping", func(ctx bot.Context, args []string) {
	if err := ctx.Send("ping"); err != nil {
		logrus.WithField("command", "ping").Errorf("Send message: %v", err)
	}
}))
