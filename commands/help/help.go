package help

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/fvdveen/bf"
	"github.com/fvdveen/mu2/commands"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("commands/help")

// Help lists all commands
var _ = commands.Register(bf.NewCommand(
	bf.Name("help"),
	bf.Trigger("help"),
	bf.Use("Lists all commands"),
	bf.Action(func(ctx bf.Context) {
		embedItems := []*discordgo.MessageEmbedField{}
		for _, com := range ctx.Bot.Commands {
			if com.Disabled() {
				continue
			}
			embedItems = append(embedItems, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("%s - %s", com.Trigger(), com.Name()),
				Value:  com.Use(),
				Inline: true,
			})

		}
		embed := &discordgo.MessageEmbed{
			Fields: embedItems,
			Title:  "Commands",
		}
		if err := ctx.SendEmbed(embed); err != nil {
			log.Errorf("help: %v", err)
		}
	}),
))
