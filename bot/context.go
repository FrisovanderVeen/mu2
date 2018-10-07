package bot

import (
	"context"
	"fmt"
	"io"

	"github.com/bwmarrin/discordgo"
	"github.com/fvdveen/mu2/db"
)

// Context holds items used by commands
type Context interface {
	Send(string) error
	SendEmbed(*discordgo.MessageEmbed) error

	Play(OpusReader) error

	Channel() (*discordgo.Channel, error)
	Guild() (*discordgo.Guild, error)
	MessageCreate() *discordgo.MessageCreate
	Session() *discordgo.Session
	Bot() *bot
	Database() db.Service
	Context() context.Context
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

func (ctx *defaultContext) Bot() *bot {
	return ctx.b
}

func (ctx *defaultContext) Database() db.Service {
	return ctx.b.db
}

func (ctx *defaultContext) Context() context.Context {
	return ctx.ctx
}

func (ctx *defaultContext) Play(or OpusReader) error {
	g, err := ctx.Guild()
	if err != nil {
		return fmt.Errorf("get guild: %v", err)
	}

	for _, vs := range g.VoiceStates {
		if vs.UserID == ctx.m.Author.ID {
			vc, err := ctx.s.ChannelVoiceJoin(g.ID, vs.ChannelID, false, true)
			if err != nil {
				return fmt.Errorf("join voice channel: %v", err)
			}

			if err := vc.Speaking(true); err != nil {
				return fmt.Errorf("set speaking status: %v", err)
			}

			for {
				o, err := or.OpusFrame()
				if err == io.EOF {
					break
				} else if err != nil {
					if err := vc.Speaking(false); err != nil {
						return fmt.Errorf("set speaking status: %v", err)
					}
					return fmt.Errorf("read opus frame: %v", err)
				}

				vc.OpusSend <- o
			}

			if err := vc.Speaking(false); err != nil {
				return fmt.Errorf("set speaking status: %v", err)
			}

			if err := vc.Disconnect(); err != nil {
				return fmt.Errorf("disconnect from voice channel: %v", err)
			}
		}
	}

	return nil
}
