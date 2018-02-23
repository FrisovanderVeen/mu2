package botFramework

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/bwmarrin/discordgo"
)

// Bot is a wrapper for a discordgo session
type Bot struct {
	Session *discordgo.Session

	ErrWriter io.Writer
	ErrPrefix func() string

	Prefix string
	Token  string

	Commands map[string]*Command
}

// NewBot creates a new bot
func NewBot(options ...OptionFunc) (*Bot, error) {
	bot := &Bot{
		Commands:  make(map[string]*Command),
		ErrWriter: os.Stderr,
		ErrPrefix: func() string {
			return ""
		},
	}

	for _, opt := range options {
		if err := opt(bot); err != nil {
			return nil, fmt.Errorf("Option: %v", err)
		}
	}

	if bot.Token == "" {
		return nil, errors.New("no token given")
	}

	sess, err := discordgo.New("Bot " + bot.Token)
	if err != nil {
		return nil, fmt.Errorf("Session: %v", err)
	}
	sess.AddHandler(bot.handleCommands)
	bot.Session = sess

	return bot, nil
}

// Error writes the error to the ErrWriter
func (b *Bot) Error(err error) {
	fmt.Fprintln(b.ErrWriter, b.ErrPrefix(), err)
}

// Open opens the discord session
func (b *Bot) Open() error {
	return b.Session.Open()
}

// Close closes the discord session
func (b *Bot) Close() error {
	return b.Session.Close()
}

// Restart closes and opens the discord session
func (b *Bot) Restart() error {
	if err := b.Session.Close(); err != nil {
		return err
	}
	return b.Session.Open()
}

// AddHandler adds a handler to the discordgo session
func (b *Bot) AddHandler(handler interface{}) {
	b.Session.AddHandler(handler)
}

// AddCommand adds a command to the bot
func (b *Bot) AddCommand(coms ...*Command) error {
	double := false

	for _, com := range coms {
		if _, ok := b.Commands[com.Trigger]; ok {
			double = true
			continue
		}
		b.Commands[com.Trigger] = com
	}

	if double {
		return ErrDoubleCommand
	}

	return nil
}

// UpdateStatus sets the bot's status if game == "" then set status to active and not playing anything
func (b *Bot) UpdateStatus(game string) error {
	return b.Session.UpdateStatus(0, game)
}
