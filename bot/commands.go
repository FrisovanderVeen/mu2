package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fvdveen/mu2/db"
	"github.com/sirupsen/logrus"
)

// Command is a command in the bot
type Command interface {
	Name() string
	Help() string
	Run(Context, []string)
}

func (b *Bot) commandHandler() func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if !strings.HasPrefix(m.Content, b.conf.Prefix) {
			return
		}

		msg := strings.Split(strings.TrimPrefix(m.Content, b.conf.Prefix), " ")

		c, err := b.Command(msg[0])
		if err != nil {
			c, err = b.dbCommand(s, m, msg[0])
			if err != nil {
				logrus.WithField("handler", "command").Errorf("Get command: %v", err)
				return
			}
		}

		ctx := &defaultContext{
			s: s,
			m: m,
			b: b,
		}

		c.Run(ctx, msg[1:])
	}
}

// NewCommand creates a new command
func NewCommand(name string, description string, action func(Context, []string)) Command {
	return &defaultCommand{
		n: name,
		d: description,
		a: action,
	}
}

type defaultCommand struct {
	n string
	d string
	a func(Context, []string)
}

func (c defaultCommand) Name() string {
	return c.n
}

func (c defaultCommand) Help() string {
	return c.d
}

func (c defaultCommand) Run(ctx Context, args []string) {
	c.a(ctx, args)
}

// HelpCommand returns the default help command
func (b *Bot) HelpCommand() Command {
	return NewCommand("help", "Sends an help message", func(c Context, _ []string) {
		var msg string

		for _, c := range b.Commands() {
			msg = fmt.Sprintf("%s`%s%s` %s\n", msg, b.conf.Prefix, c.Name(), c.Help())
		}

		if err := c.Send(msg); err != nil {
			logrus.WithField("command", "help").Errorf("Send message: %v", err)
		}
	})
}

func dbCommand(i *db.Item) Command {
	return NewCommand("", "", func(c Context, _ []string) {
		if err := c.Send(i.Response); err != nil {
			logrus.WithFields(map[string]interface{}{
				"command": "_learnable",
				"item":    i.Message,
				"guildID": i.GuildID,
			}).Errorf("Send message: %v", err)
		}
	})
}

func (b *Bot) learnCommands() []Command {
	return []Command{
		NewCommand("learn", "Teach a message-response command to the bot", func(c Context, args []string) {
			g, err := c.Guild()
			if err != nil {
				logrus.WithField("command", "learn").Errorf("Get guild: %v", err)
			}

			_, err = b.db.Get(g.ID, args[0])
			if err != nil && err != db.ErrItemNotFound {
				logrus.WithField("command", "learn").Errorf("Get item: %v", err)
				return
			} else if err == nil {
				if err := c.Send("haha, no"); err != nil {
					logrus.WithField("command", "learn").Errorf("Send message: %v", "haha, no")
					return
				}
			}

			err = b.db.New(&db.Item{
				Message:  args[0],
				Response: strings.Join(args[1:], " "),
				GuildID:  g.ID,
			})
			if err != nil {
				logrus.WithField("command", "learn").Errorf("Store item: %v", err)
			}

			if err := c.Send(fmt.Sprintf("Learned %s succesfully", args[0])); err != nil {
				logrus.WithField("command", "learn").Errorf("Send message: %v", err)
			}
		}),
	}
}

func (b *Bot) dbCommand(s *discordgo.Session, m *discordgo.MessageCreate, msg string) (Command, error) {
	ch, err := s.State.Channel(m.ChannelID)
	if err != nil {
		ch, err = s.Channel(m.ChannelID)
		if err != nil {
			return nil, err
		}
	}

	i, err := b.db.Get(ch.GuildID, msg)
	if err == db.ErrItemNotFound {
		return nil, db.ErrItemNotFound
	} else if err != nil {
		return nil, err
	}

	return dbCommand(i), nil
}