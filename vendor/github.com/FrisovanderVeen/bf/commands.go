package botFramework

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Command is a simple implementation of the CommandInterface interface
type Command struct {
	A func(Context)
	T string
	U string
	N string
	D bool
}

// NewCommand creates a new command with the given options
func NewCommand(options ...func(*Command)) *Command {
	c := &Command{}
	for _, o := range options {
		o(c)
	}

	return c
}

// Action sets the action of the new command
func Action(action func(Context)) func(*Command) {
	return func(c *Command) {
		c.A = action
	}
}

// Name sets the name of the new command
func Name(name string) func(*Command) {
	return func(c *Command) {
		c.N = name
	}
}

// Use sets the use of the new command
func Use(use string) func(*Command) {
	return func(c *Command) {
		c.U = use
	}
}

// Trigger sets the trigger of the new command
func Trigger(trigger string) func(*Command) {
	return func(c *Command) {
		c.T = trigger
	}
}

// Disabled sets if the new command is disabled
func Disabled(disabled bool) func(*Command) {
	return func(c *Command) {
		c.D = disabled
	}
}

// Action calls the command'ss action
func (c *Command) Action(ctx Context) {
	c.A(ctx)
}

// Trigger returns the command's trigger
func (c *Command) Trigger() string {
	return c.T
}

// Use returns the command's use
func (c *Command) Use() string {
	return c.U
}

// Name returns the command's name
func (c *Command) Name() string {
	return c.N
}

// Disabled returns whether the command is disabled
func (c *Command) Disabled() bool {
	return c.D
}

// CommandInterface describes a basic command
type CommandInterface interface {
	Action(Context)
	Disabled() bool
	Name() string
	Use() string
	Trigger() string
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
