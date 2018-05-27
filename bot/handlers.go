package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func (b *Bot) initHandlers() {
	b.sess.AddHandler(b.messageHandler())
	b.sess.AddHandler(b.readyHandler())
	err := loadSound()
	if err != nil {
		logrus.Error(err)
	}
}

func (b *Bot) messageHandler() func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if s.State.User.ID == m.Author.ID {
			return
		}
		if !strings.HasPrefix(m.Content, b.conf.Prefix) {
			return
		}
		msg := strings.TrimPrefix(m.Content, b.conf.Prefix)
		for _, com := range b.commands {
			if msg == com.name || strings.HasPrefix(msg, com.name) {
				msg = strings.TrimSpace(strings.TrimPrefix(com.name, msg))
				com.run(b, m, msg)
				return
			}
		}
	}
}

func (b *Bot) readyHandler() func(*discordgo.Session, *discordgo.Ready) {
	return func(s *discordgo.Session, _ *discordgo.Ready) {
		s.UpdateStatus(0, fmt.Sprintf("%shelp", b.conf.Prefix))
	}
}

func (b *Bot) sendHelp(channelID string) {
	e := NewEmbed().
		SetTitle("Mu2").
		SetDescription("Help")
	for _, com := range b.commands {
		e.AddField(com.name, com.description)
	}

	b.sess.ChannelMessageSendEmbed(channelID, e.MessageEmbed)
}

func (b *Bot) sendInfo(channelID string) {
	e := NewEmbed().
		SetTitle("Mu2").
		SetDescription("Info").
		AddField("Author", "CreepyGuy").
		AddField("Server count", strconv.Itoa(len(b.sess.State.Guilds))).
		AddField("Invite link", b.conf.InviteLink).
		AddField("Github", "https://github.com/fvdveen/mu2")
	b.sess.ChannelMessageSendEmbed(channelID, e.MessageEmbed)
}
