package bot

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/fvdveen/mu2/config"
	"github.com/fvdveen/mu2/db"
)

// Bot is a discord bot
type Bot struct {
	sess *discordgo.Session
	db   db.Service
	conf config.Config

	comMu    sync.RWMutex
	commands map[string]*Command
}

// New creates a new discord bot
func New(opts ...OptionFunc) (*Bot, error) {
	b := &Bot{
		commands: map[string]*Command{},
	}

	for _, opt := range opts {
		opt(b)
	}

	sess, err := discordgo.New("Bot " + b.conf.Discord.Token)
	if err != nil {
		return nil, err
	}
	b.sess = sess

	b.Init()

	return b, err
}

// Run runs the bot
func (b *Bot) Run() error {
	return b.sess.Open()
}

// Close closes the bot
func (b *Bot) Close() error {
	return b.sess.Close()
}

// Init initialized the bot
func (b *Bot) Init() {
	b.sess.AddHandler(b.CommandHandler)
	for _, c := range commands {
		b.AddCommand(c)
	}
}

// AddCommand adds a command to the bot
func (b *Bot) AddCommand(c *Command) func() {
	b.comMu.Lock()
	defer b.comMu.Unlock()
	b.commands[c.Name] = c

	n := c.Name

	return func() {
		b.comMu.Lock()
		defer b.comMu.Unlock()
		delete(b.commands, n)
	}
}

// GetCommand returns the command for the given trigger it will return nil if no command is found
func (b *Bot) GetCommand(n string) *Command {
	b.comMu.RLock()
	defer b.comMu.RUnlock()
	return b.commands[n]
}
