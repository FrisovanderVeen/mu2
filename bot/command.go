package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fvdveen/mu2/common"
)

var cogs = map[string]Cog{}

// RegisterCog registers a cog with the bot
// only works if RegisterCog is called before the bot is created
func RegisterCog(cog Cog) Cog {
	if _, ok := cogs[cog.Help().Trigger]; ok {
		common.GetLogger().Warnf("Cannot register multiple cogs with same trigger, trigger: %s", cog.Help().Trigger)
		return nil
	}

	for _, cat := range CogCategories {
		if cat.Name == cog.Help().Trigger {
			common.GetLogger().Warnf("Cannot register cog with same trigger as a command category, trigger: %s", cog.Help().Trigger)
			return nil
		}
	}

	cogs[cog.Help().Trigger] = cog

	return cog
}

type CogCategory struct {
	Name        string
	Description string
}

var (
	// CogCategories are all cog categories
	CogCategories = []*CogCategory{
		CogCategoryInfo,
		CogCategoryMusic,
	}
	CogCategoryInfo = &CogCategory{
		Name:        "info",
		Description: "Helpful commands on how to use the bot",
	}
	CogCategoryMusic = &CogCategory{
		Name:        "music",
		Description: "Commands related to playing music and manipulating the queue",
	}
)

type CommandHelp struct {
	Name        string
	Trigger     string
	Description string
	Module      bool
}

// Cog is a command in the bot
type Cog interface {
	Help() CommandHelp
	CogCategory() *CogCategory

	// if subcogs returns nil it means that there arent any subcommands
	SubCogs() map[string]Cog
	RunFunc() func(*Context, []string) error
}

type Module struct {
	Name        string
	Trigger     string
	Description string
	Category    *CogCategory
	Run         func(*Context, []string) error

	Sub map[string]Cog
}

func (m *Module) Help() CommandHelp {
	return CommandHelp{
		Name:        m.Name,
		Trigger:     m.Trigger,
		Description: m.Description,
		Module:      true,
	}
}

func (m *Module) CogCategory() *CogCategory {
	return m.Category
}

func (m *Module) SubCogs() map[string]Cog {
	return nil
}

func (m *Module) RunFunc() func(*Context, []string) error {
	return func(ctx *Context, args []string) error {
		if len(args) == 0 {
			return m.Run(ctx, args)
		}
		c, ok := m.SubCogs()[args[0]]
		if !ok {
			return m.Run(ctx, args)
		}
		return c.RunFunc()(ctx, args)
	}
}

type Command struct {
	Name        string
	Trigger     string
	Description string
	Category    *CogCategory
	Run         func(*Context, []string) error
}

func (c *Command) Help() CommandHelp {
	return CommandHelp{
		Name:        c.Name,
		Trigger:     c.Trigger,
		Description: c.Description,
		Module:      false,
	}
}

func (c *Command) CogCategory() *CogCategory {
	return c.Category
}

func (c *Command) SubCogs() map[string]Cog {
	return nil
}

func (c *Command) RunFunc() func(*Context, []string) error {
	return c.Run
}

type Context struct {
	Args []string

	Session       *discordgo.Session
	MessageCreate *discordgo.MessageCreate

	Bot *Bot
}

func (ctx *Context) RespondEmbed(e *discordgo.MessageEmbed) error {
	_, err := ctx.Session.ChannelMessageSendEmbed(ctx.MessageCreate.ChannelID, e)
	return err
}
