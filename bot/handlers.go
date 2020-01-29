package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fvdveen/mu2/common"
)

func (b *Bot) readyHandler() func(s *discordgo.Session, m *discordgo.Ready) {
	return func(s *discordgo.Session, m *discordgo.Ready) {
		s.UpdateStatus(0, common.GetPrefix()+"help")
	}
}

func (b *Bot) messageHandler() func(s *discordgo.Session, m *discordgo.MessageCreate) {
	log := common.GetHandlerLogger("message-create")
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if !strings.HasPrefix(m.Content, b.conf.Bot.CommandPrefix) {
			return
		}

		log.WithFields(map[string]interface{}{"guild-id": m.GuildID, "channel-id": m.ChannelID, "command": m.Content}).Debugf("Command recieved")

		ctx := b.getCommandContext(m)

		cog, ok := b.Cog(ctx.Args[0])
		if !ok {
			log.WithFields(map[string]interface{}{"guild-id": m.GuildID, "channel-id": m.ChannelID, "command": ctx.Args[0]}).Warn("Command not found")
			return
		}

		var runFunc func(*Context, []string) error
		runFunc, ctx.Args = b.getRunFunc(cog, ctx.Args[1:])
		if runFunc == nil {
			log.WithField("args", ctx.Args).Warnf("command %s has no run func", ctx.Args[0])
			return
		}

		if err := runFunc(ctx, ctx.Args); err != nil {
			common.GetHandlerLogger("message-create").WithField("cog", cog.Help().Name).Error(err)
			return
		}
	}
}
