package commands

import (
	"fmt"
	"strconv"

	bf "github.com/FrisovanderVeen/bf"
	"github.com/bwmarrin/discordgo"
)

// Info gives info about the bot
var Info = &bf.Command{
	Name:    "info",
	Trigger: "info",
	Use:     "Gives info about the bot",
	Action: func(ctx bf.Context) {
		embedItems := []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Author",
				Value:  "CreepyGuy",
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Server count",
				Value:  strconv.Itoa(len(ctx.Session.State.Guilds)),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Invite",
				Value:  "https://discordapp.com/oauth2/authorize?client_id=416569570703310850&scope=bot",
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "GitHub",
				Value:  "https://discordapp.com/oauth2/authorize?client_id=416569570703310850&scope=bot",
				Inline: true,
			},
		}
		embed := &discordgo.MessageEmbed{
			Fields:      embedItems,
			Title:       "Mu2",
			Description: "A Discord bot",
		}
		if err := ctx.SendEmbed(embed); err != nil {
			ctx.Error(fmt.Errorf("help: %v", err))
		}
	},
}
