package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fvdveen/mu2/db"
	"github.com/sirupsen/logrus"
)

// Context holds all the values required by commands
type Context struct {
	B    *Bot
	S    *discordgo.Session
	M    *discordgo.MessageCreate
	G    *discordgo.Guild
	C    *discordgo.Channel
	Args []string
}

// Command is a command used by the bot
type Command struct {
	Name   string
	Use    string
	Action func(*Context)
}

// CommandHandler is the handler the bot uses for commands
func (b *Bot) CommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(strings.ToLower(m.Content), "ligma") {
		s.ChannelMessageSend(m.ChannelID, "Ligma fucking balls lmao :joy:")
	}

	if !strings.HasPrefix(m.Content, b.conf.Discord.Prefix) {
		return
	}

	args := strings.Split(strings.TrimPrefix(m.Content, b.conf.Discord.Prefix), " ")
	c, err := s.Channel(m.ChannelID)
	if err != nil {
		logrus.Errorf("Could not get channel: %v", err)
		return
	}

	com := b.GetCommand(args[0])
	if com == nil {
		dbC, err := b.db.Command(c.GuildID, args[0])
		if err == db.ErrNoCommand {
			return
		} else if err != nil {
			logrus.Errorf("Could not get command: %s: %v", args[0], err)
			return
		}
		s.ChannelMessageSend(m.ChannelID, dbC.Response)
		return
	}
	g, err := s.Guild(c.GuildID)
	if err != nil {
		logrus.Error("Could not get guild: %v", err)
		return
	}
	ctx := &Context{
		B:    b,
		S:    s,
		M:    m,
		C:    c,
		G:    g,
		Args: args[1:],
	}
	com.Action(ctx)
}

var commands = []*Command{
	{
		Name: "learn",
		Use:  "Learns a new message-response style command",
		Action: func(c *Context) {
			if len(c.Args) < 1 {
				return
			}
			if com := c.B.GetCommand(c.Args[0]); com != nil {
				c.S.ChannelMessageSend(c.M.ChannelID, fmt.Sprintf("I'm sorry %s, I'm afraid I can't let you do that", c.M.Author.Username))
				return
			}
			_, err := c.B.db.Command(c.G.ID, c.Args[0])
			if err == nil {
				c.S.ChannelMessageSend(c.M.ChannelID, fmt.Sprintf("I'm sorry %s, I'm afraid I can't let you do that", c.M.Author.Username))
				return

			} else if err != nil && err != db.ErrNoCommand {
				logrus.Errorf("[learn] Could not get commands: %v", err)
				return
			}

			if len(c.Args) < 2 {
				return
			}

			com := &db.Command{
				GID:      c.G.ID,
				Name:     c.Args[0],
				Response: strings.Join(c.Args[1:], " "),
			}

			if err := c.B.db.AddCommand(com); err != nil {
				logrus.Errorf("[learn] Could not add command: %v", err)
				return
			}

			c.S.ChannelMessageSend(c.M.ChannelID, fmt.Sprintf("Learned command: %s succesfully", c.Args[0]))
		},
	},
	{
		Name: "unlearn",
		Use:  "Removes a learned command",
		Action: func(c *Context) {
			if len(c.Args) < 1 {
				return
			}
			if com := c.B.GetCommand(c.Args[0]); com != nil {
				c.S.ChannelMessageSend(c.M.ChannelID, "Haha, no")
				return
			}

			if err := c.B.db.RemoveCommand(c.G.ID, c.Args[0]); err != nil {
				logrus.Error("[unlearn] Could not remove command: %v", err)
				return
			}

			c.S.ChannelMessageSend(c.M.ChannelID, fmt.Sprintf("Unlearned command: %s succesfully", c.Args[0]))
		},
	},
	{
		Name: "learn-server",
		Use:  "Learns a new message-response style command to the given server-ID",
		Action: func(c *Context) {
			if len(c.Args) < 2 {
				return
			}
			if com := c.B.GetCommand(c.Args[0]); com != nil {
				c.S.ChannelMessageSend(c.M.ChannelID, fmt.Sprintf("I'm sorry %s, I'm afraid I can't let you do that", c.M.Author.Username))
				return
			}
			_, err := c.B.db.Command(c.G.ID, c.Args[0])
			if err == nil {
				c.S.ChannelMessageSend(c.M.ChannelID, fmt.Sprintf("I'm sorry %s, I'm afraid I can't let you do that", c.M.Author.Username))
				return

			} else if err != nil && err != db.ErrNoCommand {
				logrus.Errorf("[learn-server] Could not get commands: %v", err)
				return
			}

			g, err := c.S.Guild(c.Args[1])
			if err != nil {
				logrus.Errorf("[learn-server] Could not get guild: %v", err)
				return
			}

			var in bool

			for _, m := range g.Members {
				if m.User.ID == c.M.Author.ID {
					in = true
				}
			}

			if !in {
				c.S.ChannelMessageSend(c.M.ChannelID, "One of us isn't in that server")
				return
			}

			if len(c.Args) < 3 {
				return
			}

			com := &db.Command{
				GID:      c.Args[1],
				Name:     c.Args[0],
				Response: strings.Join(c.Args[2:], " "),
			}

			if err := c.B.db.AddCommand(com); err != nil {
				logrus.Errorf("[learn-server] Could not add command: %v", err)
				return
			}

			c.S.ChannelMessageSend(c.M.ChannelID, fmt.Sprintf("Learned command: %s succesfully", c.Args[0]))
		},
	},
}
