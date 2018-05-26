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
		switch {
		case msg == "info" || strings.HasPrefix(msg, "info "):
			b.sendInfo(m.ChannelID)
		case msg == "help" || strings.HasPrefix(msg, "help "):
			b.sendHelp(m.ChannelID)
		case msg == "airhorn" || strings.HasPrefix(msg, "airhorn "):
			c, err := s.State.Channel(m.ChannelID)
			if err != nil {
				logrus.Error(err)
				s.ChannelMessageSend(m.ChannelID, "oops something went wrong, try again later")
				return
			}
			b.playAirhorn(m.Author.ID, m.ChannelID, c.GuildID)
		case msg == "skip" || strings.HasPrefix(msg, "skip "):
			c, err := s.State.Channel(m.ChannelID)
			if err != nil {
				logrus.Error(err)
				s.ChannelMessageSend(m.ChannelID, "oops something went wrong, try again later")
				return
			}
			vh, ok := b.voiceHandlers[c.GuildID]
			if !ok {
				s.ChannelMessageSend(m.ChannelID, "Bot has to playing to skip")
				return
			}
			vh.skipChan <- 0
		case msg == "stop" || strings.HasPrefix(msg, "stop "):
			c, err := s.State.Channel(m.ChannelID)
			if err != nil {
				logrus.Error(err)
				s.ChannelMessageSend(m.ChannelID, "oops something went wrong, try again later")
				return
			}
			vh, ok := b.voiceHandlers[c.GuildID]
			if !ok {
				s.ChannelMessageSend(m.ChannelID, "Bot has to playing to stop")
				return
			}
			vh.stopChan <- 0
		case msg == "loop" || strings.HasPrefix(msg, "loop "):
			c, err := s.State.Channel(m.ChannelID)
			if err != nil {
				logrus.Error(err)
				s.ChannelMessageSend(m.ChannelID, "oops something went wrong, try again later")
				return
			}
			vh, ok := b.voiceHandlers[c.GuildID]
			if !ok {
				s.ChannelMessageSend(m.ChannelID, "Bot has to playing to loop")
				return
			}
			vh.loopChan <- 0
		case msg == "play" || strings.HasPrefix(msg, "play "):
			c, err := s.State.Channel(m.ChannelID)
			if err != nil {
				logrus.Error(err)
				s.ChannelMessageSend(m.ChannelID, "oops something went wrong, try again later")
				return
			}
			vh, ok := b.voiceHandlers[c.GuildID]
			if !ok {
				s.ChannelMessageSend(m.ChannelID, "Bot has to playing to play")
				return
			}
			vh.playChan <- 0
		case msg == "pause" || strings.HasPrefix(msg, "pause "):
			c, err := s.State.Channel(m.ChannelID)
			if err != nil {
				logrus.Error(err)
				s.ChannelMessageSend(m.ChannelID, "oops something went wrong, try again later")
				return
			}
			vh, ok := b.voiceHandlers[c.GuildID]
			if !ok {
				s.ChannelMessageSend(m.ChannelID, "Bot has to playing to pause")
				return
			}
			vh.pauseChan <- 0
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
		SetDescription("Help").
		AddField("help", "Gives help on all commands").
		AddField("info", "Gives info on the bot").
		AddField("skip", "Skips the currently playing audio").
		AddField("stop", "Stops all audio and leaves the voice channel").
		AddField("loop", "Loop the queue, new elements can still be added").
		AddField("play", "Plays the currently playing audio").
		AddField("pause", "Pauses the currently playing audio").
		AddField("airhorn", "Plays an annoying airhorn sound")

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
