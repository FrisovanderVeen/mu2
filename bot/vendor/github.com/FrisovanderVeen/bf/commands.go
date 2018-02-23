package botFramework

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Command describes an action that the bot performs
type Command struct {
	Action   func(Context)
	Disabled bool
	Trigger  string

	Name string
	Use  string
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
		if com.Disabled {
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
