package info

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	humanize "github.com/dustin/go-humanize"
	"github.com/fvdveen/bf"
	"github.com/fvdveen/mu2/commands"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("commands/info")

// Info gives info about the bot
var _ = commands.Register(bf.NewCommand(
	bf.Name("info"),
	bf.Trigger("info"),
	bf.Use("Gives info about the bot"),
	bf.Action(func(ctx bf.Context) {
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
				Name:   "Running since",
				Value:  humanize.Time(ctx.Bot.Time),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Invite",
				Value:  "https://discordapp.com/oauth2/authorize?client_id=416569570703310850&scope=bot",
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "GitHub",
				Value:  "https://github.com/fvdveen/mu2",
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Version",
				Value:  commands.VERSION,
				Inline: true,
			},
		}
		embed := &discordgo.MessageEmbed{
			Fields:      embedItems,
			Title:       "Mu2",
			Description: "A Discord bot",
		}
		if err := ctx.SendEmbed(embed); err != nil {
			log.Errorf("help: %v", err)
		}
	}),
))
