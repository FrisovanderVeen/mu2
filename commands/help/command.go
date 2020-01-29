package help

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fvdveen/mu2/bot"
	"github.com/fvdveen/mu2/common"
)

var log = common.GetCommandLogger("help")
var _ = bot.RegisterCog(Command)

var Command = &bot.Command{
	Name:        "Help",
	Trigger:     "help",
	Description: "Gives the use of commands and modules",
	Category:    bot.CogCategoryInfo,
	Run: func(ctx *bot.Context, args []string) error {
		log.Debug("running command", args)
		if len(args) == 0 {
			return common.CommandError("help", sendDefaultMessage(ctx))
		}

		c := strings.TrimPrefix(args[0], common.GetPrefix())

		cog, ok := ctx.Bot.Cog(c)
		if !ok {
			cog, ok = ctx.Bot.CogByName(c)
			if !ok {
				cat, ok := getCommandCategory(c)
				if !ok {
					return common.CommandError("help", sendCommandNotFoundMessage(ctx, c))
				}
				return common.CommandError("help", sendCategoryMessage(ctx, cat))
			}
		}

		return common.CommandError("help", sendHelpMessage(ctx, cog))
	},
}

func sendDefaultMessage(ctx *bot.Context) error {
	e := &discordgo.MessageEmbed{
		Title:       "Help",
		Description: "Send " + common.GetPrefix() + "command to use a command e.g. " + common.GetPrefix() + "help",
	}

	for cat, cogs := range ctx.Bot.CogsByCategory() {
		var coms []string

		for _, cog := range cogs {
			coms = append(coms, "`"+common.GetPrefix()+cog.Help().Trigger+"`")
		}

		e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
			Name:  cat.Name,
			Value: strings.Join(coms, " "),
		})
	}

	return ctx.RespondEmbed(e)
}

func sendCommandNotFoundMessage(ctx *bot.Context, command string) error {
	e := &discordgo.MessageEmbed{
		Title:       "Unknown command",
		Description: "Command " + command + " is not recognised",
	}
	return ctx.RespondEmbed(e)
}

func sendHelpMessage(ctx *bot.Context, cog bot.Cog) error {
	help := cog.Help()
	e := &discordgo.MessageEmbed{
		Title:       help.Name,
		Description: "`" + common.GetPrefix() + help.Trigger + "` - " + help.Description,
	}
	return ctx.RespondEmbed(e)
}

func sendCategoryMessage(ctx *bot.Context, cat *bot.CogCategory) error {
	e := &discordgo.MessageEmbed{
		Title:       cat.Name,
		Description: cat.Description,
	}

	for _, cog := range ctx.Bot.CogsForCategory(cat) {
		help := cog.Help()
		e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
			Name:  help.Name,
			Value: "`" + common.GetPrefix() + help.Trigger + "` - " + help.Description,
		})
	}

	return ctx.RespondEmbed(e)
}

func getCommandCategory(com string) (*bot.CogCategory, bool) {
	for _, cat := range bot.CogCategories {
		if cat.Name == com {
			return cat, true
		}
	}

	return nil, false
}
