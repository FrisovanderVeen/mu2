package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Context holds items used by commands
type Context interface {
	Send(string) error
	SendEmbed(*discordgo.MessageEmbed) error

	Channel() (*discordgo.Channel, error)
	Guild() (*discordgo.Guild, error)
	MessageCreate() *discordgo.MessageCreate
	Session() *discordgo.Session
	Bot() *Bot
}

type defaultContext struct {
	m *discordgo.MessageCreate
	s *discordgo.Session
	b *Bot
}

func (ctx defaultContext) Send(s string) error {
	fmt.Println(s)
	_, err := ctx.s.ChannelMessageSend(ctx.m.ChannelID, s)
	return err
}

func (ctx defaultContext) SendEmbed(e *discordgo.MessageEmbed) error {
	_, err := ctx.s.ChannelMessageSendEmbed(ctx.m.ChannelID, e)
	return err
}

func (ctx defaultContext) Session() *discordgo.Session {
	return ctx.s
}

func (ctx defaultContext) MessageCreate() *discordgo.MessageCreate {
	return ctx.m
}

func (ctx defaultContext) Guild() (*discordgo.Guild, error) {
	c, err := ctx.Channel()
	if err != nil {
		return nil, fmt.Errorf("get channel: %v", err)
	}

	g, err := ctx.s.State.Guild(c.GuildID)
	if err != nil {
		g, err = ctx.s.Guild(c.GuildID)
		if err != nil {
			return nil, err
		}
	}

	return g, nil
}

func (ctx defaultContext) Channel() (*discordgo.Channel, error) {
	c, err := ctx.s.State.Channel(ctx.m.ChannelID)
	if err != nil {
		c, err = ctx.s.Channel(ctx.m.ChannelID)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (ctx defaultContext) Bot() *Bot {
	return ctx.b
}