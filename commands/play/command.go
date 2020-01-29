package play

import (
	"fmt"
	"strings"

	"github.com/fvdveen/mu2/bot"
	"github.com/fvdveen/mu2/common"
	"github.com/fvdveen/mu2/voice/buffer"
	"github.com/fvdveen/mu2/voice/youtube"
)

var log = common.GetCommandLogger("play")
var _ = bot.RegisterCog(Command)

var Command = &bot.Command{
	Name:        "Play",
	Trigger:     "play",
	Description: "",
	Category:    bot.CogCategoryMusic,
	Run: func(ctx *bot.Context, args []string) error {
		vh, err := ctx.GetVoiceHandler()
		if err != nil {
			return common.CommandError("play", fmt.Errorf("get voice handler: %v", err))
		}

		inf, err := youtube.Search(strings.Join(args, " "))
		if err != nil {
			return common.CommandError("play", fmt.Errorf("search youtube handler: %v", err))
		}

		vi, err := youtube.NewVoiceItem(inf)
		if err != nil {
			return common.CommandError("play", fmt.Errorf("create voice item: %v", err))
		}

		buf := buffer.New(vi, buffer.WithAsync())
		vh.Play(buf)

		return nil
	},
}
