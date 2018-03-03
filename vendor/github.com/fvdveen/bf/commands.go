package bf

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Command describes a basic command
type Command interface {
	Action(Context)
	Disabled() bool
	Name() string
	Use() string
	Trigger() string
}

// Com is a simple implementation of the Command interface
type Com struct {
	A func(Context)
	T string
	U string
	N string
	D bool
}

// NewCommand creates a new command with the given options
func NewCommand(options ...func(*Com)) *Com {
	c := &Com{}
	for _, o := range options {
		o(c)
	}

	return c
}

// Action sets the action of the new command
func Action(action func(Context)) func(*Com) {
	return func(c *Com) {
		c.A = action
	}
}

// Name sets the name of the new command
func Name(name string) func(*Com) {
	return func(c *Com) {
		c.N = name
	}
}

// Use sets the use of the new command
func Use(use string) func(*Com) {
	return func(c *Com) {
		c.U = use
	}
}

// Trigger sets the trigger of the new command
func Trigger(trigger string) func(*Com) {
	return func(c *Com) {
		c.T = trigger
	}
}

// Disabled sets if the new command is disabled
func Disabled(disabled bool) func(*Com) {
	return func(c *Com) {
		c.D = disabled
	}
}

// Action calls the command'ss action
func (c *Com) Action(ctx Context) {
	c.A(ctx)
}

// Trigger returns the command's trigger
func (c *Com) Trigger() string {
	return c.T
}

// Use returns the command's use
func (c *Com) Use() string {
	return c.U
}

// Name returns the command's name
func (c *Com) Name() string {
	return c.N
}

// Disabled returns whether the command is disabled
func (c *Com) Disabled() bool {
	return c.D
}

func (b *Bot) handleCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	message := m.Content
	if !strings.HasPrefix(message, b.Prefix) {
		return
	}
	message = strings.TrimPrefix(message, b.Prefix)

	for trigger, com := range b.Commands {
		if com.Disabled() {
			continue
		}
		if strings.HasPrefix(message, trigger) {
			message = strings.TrimPrefix(message, trigger)
			if strings.HasPrefix(message, " ") {
				message = strings.TrimPrefix(message, " ")
			}

			ctx := Context{
				Bot:           b,
				Session:       s,
				MessageCreate: m,
				Message:       message,
			}
			go com.Action(ctx)
		}
	}
}
