package bot

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/fvdveen/mu2/common"
	"github.com/fvdveen/mu2/config"
)

var whiteSpaceRegex = regexp.MustCompile(`\s+`)

// Bot is a discord bot
type Bot struct {
	session *discordgo.Session
	conf    config.Config

	cogs map[string]Cog

	voiceHandlers   map[string]*VoiceHandler
	voiceHandlersMu sync.RWMutex
}

// OptionFunc is a option for the bot
type OptionFunc func(*Bot)

// WithConfig sets the config for the bot
func WithConfig(conf config.Config) OptionFunc {
	return func(b *Bot) {
		b.conf = conf
	}
}

// New creates a new bot
func New(opts ...OptionFunc) (*Bot, error) {
	b := &Bot{
		cogs:          map[string]Cog{},
		voiceHandlers: map[string]*VoiceHandler{},
	}

	for _, opt := range opts {
		opt(b)
	}

	var err error
	b.session, err = discordgo.New("Bot " + b.conf.Bot.Token)
	if err != nil {
		return nil, fmt.Errorf("create session: %v", err)
	}

	b.initSession()

	b.initCogs()

	return b, nil
}

// Open opens the discord connection
func (b *Bot) Open() error {
	return b.session.Open()
}

// Close closes the discord connection and stops the bot
func (b *Bot) Close() error {
	return b.session.Close()
}

func (b *Bot) initSession() {
	b.session.AddHandler(b.readyHandler())
	b.session.AddHandler(b.messageHandler())
}

func (b *Bot) initCogs() {
	b.cogs = cogs
}

func (b *Bot) getCommandContext(m *discordgo.MessageCreate) *Context {
	ctx := &Context{
		Bot:           b,
		MessageCreate: m,
		Session:       b.session,
	}

	ctx.Args = strings.Split(
		strings.TrimPrefix(
			whiteSpaceRegex.ReplaceAllString(m.Content, " "),
			common.GetPrefix(),
		),
		" ",
	)

	return ctx
}

func (b *Bot) getRunFunc(cog Cog, args []string) (func(*Context, []string) error, []string) {
	if len(args) == 0 {
		return cog.RunFunc(), args
	}

	if subs := cog.SubCogs(); subs != nil {
		if len(args) < 1 {
			return cog.RunFunc(), args
		}

		sub := subs[args[0]]
		if sub == nil {
			return cog.RunFunc(), args
		}

		return b.getRunFunc(sub, args[1:])
	}

	return cog.RunFunc(), args
}

// Cog returns the cog for trigger com
func (b *Bot) Cog(com string) (Cog, bool) {
	c, ok := b.cogs[com]
	return c, ok
}

// Cogs returns a list of all cogs
func (b *Bot) Cogs() []Cog {
	cogs := []Cog{}

	for _, c := range b.cogs {
		cogs = append(cogs, c)
	}

	return cogs
}

// CogsByCategory returns a list of all cogs with category cat
func (b *Bot) CogsByCategory() map[*CogCategory][]Cog {
	cogs := map[*CogCategory][]Cog{}

	for _, c := range b.cogs {
		cogs[c.CogCategory()] = append(cogs[c.CogCategory()], c)
	}

	return cogs
}

// CogsForCategory returns all cogs with category cat
func (b *Bot) CogsForCategory(cat *CogCategory) []Cog {
	cogs := []Cog{}

	for _, cog := range b.cogs {
		if cog.CogCategory() == cat {
			cogs = append(cogs, cog)
		}
	}

	return cogs
}

// CogByName returns the cog with name name
func (b *Bot) CogByName(name string) (Cog, bool) {
	for _, c := range b.cogs {
		if c.Help().Name == name {
			return c, true
		}
	}

	return nil, false
}
