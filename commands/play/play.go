package play

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fvdveen/mu2/bot"
	"github.com/fvdveen/mu2/services/encode"
	"github.com/fvdveen/mu2/services/search"
	"github.com/sirupsen/logrus"
)

type command struct {
	ss search.Service
	es encode.Service
}

func New(ss search.Service, es encode.Service) bot.Command {
	return &command{
		ss: ss,
		es: es,
	}
}

func (c *command) Name() string {
	return "play"
}

func (c *command) Help() string {
	return "plays stuff"
}

func (c *command) Run(ctx bot.Context, args []string) {
	v, err := c.ss.Search(ctx.Context(), strings.Join(args, " "))
	if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "play"}).Errorf("Search: %v", err)
		return
	}

	or, err := c.es.Encode(ctx.Context(), v.URL)
	if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "play"}).Errorf("Encode: %v", err)
		return
	}

	if err := ctx.Play(bot.NewVideo(v, or, ctx)); err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "play"}).Errorf("Play: %v", err)
		return
	}

	e := &discordgo.MessageEmbed{
		Description: fmt.Sprintf("**Queued** [%s](%s)\nBy: %s", v.Name, v.URL, v.Author),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: v.ThumbnailURL,
		},
	}

	if err := ctx.SendEmbed(e); err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "play"}).Errorf("Send embed: %v", err)
		return
	}
}
