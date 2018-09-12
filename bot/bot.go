package bot

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/fvdveen/mu2/config"
	"github.com/fvdveen/mu2/db"
	"github.com/sirupsen/logrus"
)

// OptionFunc sets an option in the bot
type OptionFunc func(*Bot)

// WithConfig sets the bots config
func WithConfig(conf config.Discord) OptionFunc {
	return func(b *Bot) {
		b.conf = conf
	}
}

// WithDB sets the bots db
func WithDB(db db.Service) OptionFunc {
	return func(b *Bot) {
		b.db = db
	}
}

// Bot is a discord bot
type Bot struct {
	sess *discordgo.Session
	conf config.Discord

	db db.Service

	commu    sync.RWMutex
	commands map[string]Command
}

// New creates a bot
func New(ops ...OptionFunc) (*Bot, error) {
	b := &Bot{
		commands: make(map[string]Command),
	}

	for _, o := range ops {
		o(b)
	}

	sess, err := discordgo.New("Bot " + b.conf.Token)
	if err != nil {
		return nil, fmt.Errorf("opening discord session: %v", err)
	}
	b.sess = sess

	b.init()

	return b, nil
}

// Open opens the session
func (b *Bot) Open() error {
	return b.sess.Open()
}

// Close closes the session
func (b *Bot) Close() error {
	return b.sess.Close()
}

// AddCommand adds a command to the bot
func (b *Bot) AddCommand(cs ...Command) error {
	b.commu.Lock()
	defer b.commu.Unlock()

	for _, c := range cs {
		if _, ok := b.commands[c.Name()]; ok {
			return fmt.Errorf("command registered twice: %s", c.Name())
		}
		b.commands[c.Name()] = c
	}

	return nil
}

// Command returns the command with the given trigger
func (b *Bot) Command(n string) (Command, error) {
	b.commu.RLock()
	defer b.commu.RUnlock()

	c, ok := b.commands[n]
	if !ok {
		return nil, fmt.Errorf("command not found: %s", n)
	}

	return c, nil
}

// Commands returns all commands
func (b *Bot) Commands() []Command {
	b.commu.RLock()
	defer b.commu.RUnlock()

	cs := []Command{}
	for _, c := range b.commands {
		cs = append(cs, c)
	}

	return cs
}

func (b *Bot) readyHandler() func(*discordgo.Session, *discordgo.Ready) {
	return func(s *discordgo.Session, _ *discordgo.Ready) {
		if err := s.UpdateStatus(0, fmt.Sprintf("%shelp", b.conf.Prefix)); err != nil {
			logrus.WithField("handler", "ready").Errorf("Update status: %v", err)
		}
	}
}

func (b *Bot) init() {
	b.sess.AddHandler(b.commandHandler())
	b.sess.AddHandler(b.readyHandler())

	b.AddCommand(b.HelpCommand())
	b.AddCommand(b.learnCommands()...)
}
