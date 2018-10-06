package play

import (
	"strings"

	"github.com/fvdveen/mu2/services/search"
	"github.com/fvdveen/mu2/bot"
	"github.com/sirupsen/logrus"
)

type command struct {
	s search.Service
}

func New(ss search.Service) bot.Command {
	return &command{
		s: ss,
	}
}

func (c *command) Name() string {
	return "play"
}

func (c *command) Help() string {
	return "plays stuff"
}

func (c *command) Run(ctx bot.Context, args []string) {
	v, err := c.s.Search(ctx.Context(), strings.Join(args, " "))
	if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "play"}).Errorf("Search: %v", err)
		return
	}

	if err := ctx.Send(v.URL); err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "play"}).Errorf("Send message: %v", err)
		return
	}
}