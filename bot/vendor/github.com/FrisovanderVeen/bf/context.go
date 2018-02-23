package botFramework

import (
	"github.com/bwmarrin/discordgo"
)

// Context is a struct with data for commands
type Context struct {
	Bot           *Bot
	Session       *discordgo.Session
	MessageCreate *discordgo.MessageCreate

	Message string
}

// SendMessage sends a message to the channel that triggerd the command
func (ctx *Context) SendMessage(message string) error {
	if _, err := ctx.Session.ChannelMessageSend(ctx.MessageCreate.ChannelID, message); err != nil {
		return err
	}
	return nil
}

// SendEmbed sends a discordgo embed to the channel that triggerd the command
func (ctx *Context) SendEmbed(embed *discordgo.MessageEmbed) error {
	if _, err := ctx.Session.ChannelMessageSendEmbed(ctx.MessageCreate.ChannelID, embed); err != nil {
		return err
	}
	return nil
}

// Error writes a error to the bot
func (ctx *Context) Error(err error) {
	ctx.Bot.Error(err)
}

// GetVoiceState tries to get the voice state of the person who triggered the command
func (ctx *Context) GetVoiceState() (*discordgo.VoiceState, error) {
	for _, guild := range ctx.Session.State.Guilds {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == ctx.MessageCreate.Author.ID {
				return vs, nil
			}
		}
	}
	return nil, ErrVSNotFound
}

// GetVoiceConn joins the voice channel of the given voice state
func (ctx *Context) GetVoiceConn(vs *discordgo.VoiceState, mute bool, deaf bool) (*discordgo.VoiceConnection, error) {
	return ctx.Session.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, mute, deaf)
}

// JoinVoiceChannel joins the voice channel of the person who triggered the command
func (ctx *Context) JoinVoiceChannel(mute bool, deaf bool) (*discordgo.VoiceConnection, error) {
	vs, err := ctx.GetVoiceState()
	if err != nil {
		return nil, err
	}

	return ctx.GetVoiceConn(vs, mute, deaf)
}
