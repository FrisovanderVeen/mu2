package bot

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/fvdveen/mu2/config"
	"github.com/fvdveen/mu2/db"
	"github.com/sirupsen/logrus"
)

// Bot is a discord bot
type Bot interface {
	// Open the session
	Open() error
	// Close the session
	Close() error
	AddCommand(...Command) error
	Command(string) (Command, error)
	Commands() []Command
	RemoveCommand(string) error
	SetPrefix(string) error
	Prefix() string
	SetToken(string) error
	// Check if the bot is alive
	Ping() error
}

// OptionFunc sets an option in the bot
type OptionFunc func(*bot)

// WithConfig sets the bots config
func WithConfig(conf config.Bot) OptionFunc {
	return func(b *bot) {
		b.conf = conf
	}
}

// WithDB sets the bots db
func WithDB(db db.Service) OptionFunc {
	return func(b *bot) {
		b.db = db
	}
}

// bot implements Bot
type bot struct {
	sess   *discordgo.Session
	sessMu sync.RWMutex
	conf   config.Bot
	confMu sync.RWMutex

	db db.Service

	commu    sync.RWMutex
	commands map[string]Command
}

// New creates a bot
func New(ops ...OptionFunc) (Bot, error) {
	b := &bot{
		commands: make(map[string]Command),
	}

	for _, o := range ops {
		o(b)
	}

	sess, err := discordgo.New("Bot " + b.conf.Discord.Token)
	if err != nil {
		return nil, fmt.Errorf("opening discord session: %v", err)
	}
	b.sess = sess

	b.init()

	b.AddCommand(b.HelpCommand())
	b.AddCommand(b.learnCommands()...)
	b.AddCommand(b.InfoCommand())

	return b, nil
}

func (b *bot) init() {
	b.sess.AddHandler(b.commandHandler())
	b.sess.AddHandler(b.readyHandler())
}

func (b *bot) readyHandler() func(*discordgo.Session, *discordgo.Ready) {
	return func(s *discordgo.Session, _ *discordgo.Ready) {
		b.confMu.RLock()
		defer b.confMu.RUnlock()
		if err := s.UpdateStatus(0, fmt.Sprintf("%shelp", b.conf.Prefix)); err != nil {
			logrus.WithFields(map[string]interface{}{"type": "handler", "handler": "ready"}).Errorf("Update status: %v", err)
		}
	}
}

// Open opens the session
func (b *bot) Open() error {
	b.sessMu.RLock()
	defer b.sessMu.RUnlock()
	return b.sess.Open()
}

// Close closes the session
func (b *bot) Close() error {
	b.sessMu.RLock()
	defer b.sessMu.RUnlock()
	return b.sess.Close()
}

// AddCommand adds a command to the bot
func (b *bot) AddCommand(cs ...Command) error {
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
func (b *bot) Command(n string) (Command, error) {
	b.commu.RLock()
	defer b.commu.RUnlock()

	c, ok := b.commands[n]
	if !ok {
		return nil, fmt.Errorf("command not found: %s", n)
	}

	return c, nil
}

// Commands returns all commands
func (b *bot) Commands() []Command {
	b.commu.RLock()
	defer b.commu.RUnlock()

	cs := []Command{}
	for _, c := range b.commands {
		cs = append(cs, c)
	}

	return cs
}

func (b *bot) RemoveCommand(n string) error {
	switch n {
	case "learn", "help", "unlearn", "info":
		return fmt.Errorf("cant unlearn default command: %s", n)
	}

	b.commu.Lock()
	defer b.commu.Unlock()

	delete(b.commands, n)
	return nil
}

func (b *bot) Ping() error {
	b.sessMu.RLock()
	_, err := b.sess.User(b.sess.State.User.ID)
	b.sessMu.RUnlock()
	return err
}

func (b *bot) SetPrefix(p string) error {
	b.sessMu.RLock()
	if err := b.sess.UpdateStatus(0, fmt.Sprintf("%shelp", p)); err != nil {
		return err
	}
	b.sessMu.RUnlock()

	b.confMu.Lock()
	b.conf.Prefix = p
	b.confMu.Unlock()

	return nil
}

func (b *bot) Prefix() string {
	b.confMu.Lock()
	defer b.confMu.Unlock()
	return b.conf.Prefix
}

func (b *bot) SetToken(t string) error {
	b.sessMu.Lock()
	defer b.sessMu.Unlock()
	if err := b.sess.Close(); err != nil {
		return fmt.Errorf("close session: %v", err)
	}

	b.confMu.Lock()
	b.conf.Discord.Token = t
	b.confMu.Unlock()

	sess, err := discordgo.New("Bot " + b.conf.Discord.Token)
	if err != nil {
		return fmt.Errorf("opening discord session: %v", err)
	}
	b.sess = sess

	b.init()

	if err := b.sess.Open(); err != nil {
		return fmt.Errorf("open session: %v", err)
	}

	return nil
}
