package commands

import (
	"fmt"

	bf "github.com/FrisovanderVeen/bf"
	"github.com/bwmarrin/discordgo"
)

// Help lists all commands
var Help = &bf.Command{
	Name:    "help",
	Trigger: "help",
	Use:     "Lists all commands",
	Action: func(ctx bf.Context) {
		embedItems := []*discordgo.MessageEmbedField{}
		for _, com := range ctx.Bot.Commands {
			if com.Disabled {
				continue
			}
			embedItems = append(embedItems, &discordgo.MessageEmbedField{
				Name:   com.Name,
				Value:  com.Use,
				Inline: true,
			})

		}
		embed := &discordgo.MessageEmbed{
			Fields: embedItems,
			Title:  "Commands",
		}
		if err := ctx.SendEmbed(embed); err != nil {
			ctx.Error(fmt.Errorf("help: %v", err))
		}
	},
}
