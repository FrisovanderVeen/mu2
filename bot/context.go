package bot

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/fvdveen/mu2/db"
)

var (
	// ErrVoiceStateNotFound is used when a voice state is not found
	ErrVoiceStateNotFound = errors.New("voice state not found")
)

// Context holds items used by commands
type Context interface {
	Send(string) error
	SendEmbed(*discordgo.MessageEmbed) error

	Play(Video) error
	VoiceHandler() (VoiceHandler, error)

	Channel() (*discordgo.Channel, error)
	Guild() (*discordgo.Guild, error)
	MessageCreate() *discordgo.MessageCreate
	Session() *discordgo.Session
	Bot() Bot
	Database() db.Service
	Context() context.Context
}

// NewContext creates a new Context
func (b *bot) NewContext(ctx context.Context, m *discordgo.MessageCreate, s *discordgo.Session) Context {
	return &defaultContext{
		s:   s,
		m:   m,
		b:   b,
		ctx: ctx,
	}
}

type defaultContext struct {
	m   *discordgo.MessageCreate
	s   *discordgo.Session
	b   *bot
	ctx context.Context
}

func (ctx *defaultContext) Send(s string) error {
	_, err := ctx.s.ChannelMessageSend(ctx.m.ChannelID, s)
	return err
}

func (ctx *defaultContext) SendEmbed(e *discordgo.MessageEmbed) error {
	_, err := ctx.s.ChannelMessageSendEmbed(ctx.m.ChannelID, e)
	return err
}

func (ctx *defaultContext) Session() *discordgo.Session {
	return ctx.s
}

func (ctx *defaultContext) MessageCreate() *discordgo.MessageCreate {
	return ctx.m
}

func (ctx *defaultContext) Guild() (*discordgo.Guild, error) {
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

func (ctx *defaultContext) Channel() (*discordgo.Channel, error) {
	c, err := ctx.s.State.Channel(ctx.m.ChannelID)
	if err != nil {
		c, err = ctx.s.Channel(ctx.m.ChannelID)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (ctx *defaultContext) Bot() Bot {
	return ctx.b
}

func (ctx *defaultContext) Database() db.Service {
	return ctx.b.db
}

func (ctx *defaultContext) Context() context.Context {
	return ctx.ctx
}

func (ctx *defaultContext) Play(v Video) error {
	vh, err := ctx.VoiceHandler()
	if err == ErrVoiceStateNotFound {
		return ErrVoiceStateNotFound
	} else if err != nil {
		return fmt.Errorf("get voice handler: %v", err)
	}

	vh.Play(v)
	return nil
}

func (ctx *defaultContext) VoiceHandler() (VoiceHandler, error) {
	g, err := ctx.Guild()
	if err != nil {
		return nil, fmt.Errorf("get guild: %v", err)
	}

	var vs *discordgo.VoiceState
	var found = false

	for _, vs = range g.VoiceStates {
		if vs.UserID == ctx.m.Author.ID {
			found = true

			break
		}
	}

	if !found {
		return nil, ErrVoiceStateNotFound
	}

	return ctx.b.VoiceHandler(g.ID, vs.ChannelID)
}
