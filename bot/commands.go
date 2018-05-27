package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type command struct {
	name        string
	description string

	run func(*Bot, *discordgo.MessageCreate, string)
}

var commands = []*command{
	{
		name:        "help",
		description: "Gives help on all commands",
		run: func(b *Bot, m *discordgo.MessageCreate, msg string) {
			b.sendHelp(m.ChannelID)
		},
	},
	{
		name:        "info",
		description: "Gives info on the bot",
		run: func(b *Bot, m *discordgo.MessageCreate, msg string) {
			b.sendInfo(m.ChannelID)
		},
	},
	{
		name:        "airhorn",
		description: "Plays an annoying airhorn sound",
		run: func(b *Bot, m *discordgo.MessageCreate, msg string) {
			c, err := b.sess.State.Channel(m.ChannelID)
			if err != nil {
				logrus.Error(err)
				b.sess.ChannelMessageSend(m.ChannelID, "oops something went wrong, try again later")
				return
			}
			b.playAirhorn(m.Author.ID, m.ChannelID, c.GuildID)
		},
	},
	{
		name:        "skip",
		description: "Skips the currently playing audio",
		run: func(b *Bot, m *discordgo.MessageCreate, msg string) {
			c, err := b.sess.State.Channel(m.ChannelID)
			if err != nil {
				logrus.Error(err)
				b.sess.ChannelMessageSend(m.ChannelID, "oops something went wrong, try again later")
				return
			}
			b.voiceMu.Lock()
			vh, ok := b.voiceHandlers[c.GuildID]
			if !ok {
				b.voiceMu.Unlock()
				b.sess.ChannelMessageSend(m.ChannelID, "Bot has to playing to skip")
				return
			}
			vh.skipChan <- 0
			b.voiceMu.Unlock()
		},
	},
	{
		name:        "stop",
		description: "Stops all audio and leaves the voice channel",
		run: func(b *Bot, m *discordgo.MessageCreate, msg string) {
			c, err := b.sess.State.Channel(m.ChannelID)
			if err != nil {
				logrus.Error(err)
				b.sess.ChannelMessageSend(m.ChannelID, "oops something went wrong, try again later")
				return
			}
			b.voiceMu.Lock()
			vh, ok := b.voiceHandlers[c.GuildID]
			if !ok {
				b.voiceMu.Unlock()
				b.sess.ChannelMessageSend(m.ChannelID, "Bot has to playing to stop")
				return
			}
			vh.stopChan <- 0
			b.voiceMu.Unlock()
		},
	},
	{
		name:        "loop",
		description: "Loop the queue, new elements can still be added",
		run: func(b *Bot, m *discordgo.MessageCreate, msg string) {
			c, err := b.sess.State.Channel(m.ChannelID)
			if err != nil {
				logrus.Error(err)
				b.sess.ChannelMessageSend(m.ChannelID, "oops something went wrong, try again later")
				return
			}
			b.voiceMu.Lock()
			vh, ok := b.voiceHandlers[c.GuildID]
			if !ok {
				b.voiceMu.Unlock()
				b.sess.ChannelMessageSend(m.ChannelID, "Bot has to playing to loop")
				return
			}
			vh.loopChan <- 0
			b.voiceMu.Unlock()
		},
	},
	{
		name:        "play",
		description: "Plays the currently playing audio",
		run: func(b *Bot, m *discordgo.MessageCreate, msg string) {
			c, err := b.sess.State.Channel(m.ChannelID)
			if err != nil {
				logrus.Error(err)
				b.sess.ChannelMessageSend(m.ChannelID, "oops something went wrong, try again later")
				return
			}
			b.voiceMu.Lock()
			vh, ok := b.voiceHandlers[c.GuildID]
			if !ok {
				b.voiceMu.Unlock()
				b.sess.ChannelMessageSend(m.ChannelID, "Bot has to playing to play")
				return
			}
			vh.playChan <- 0
			b.voiceMu.Unlock()
		},
	},
	{
		name:        "pause",
		description: "Pauses the currently playing audio",
		run: func(b *Bot, m *discordgo.MessageCreate, msg string) {
			c, err := b.sess.State.Channel(m.ChannelID)
			if err != nil {
				logrus.Error(err)
				b.sess.ChannelMessageSend(m.ChannelID, "oops something went wrong, try again later")
				return
			}
			b.voiceMu.Lock()
			vh, ok := b.voiceHandlers[c.GuildID]
			if !ok {
				b.voiceMu.Unlock()
				b.sess.ChannelMessageSend(m.ChannelID, "Bot has to playing to pause")
				return
			}
			vh.pauseChan <- 0
			b.voiceMu.Unlock()
		},
	},
	{
		name:        "repeat",
		description: "Repeat the currently playing item",
		run: func(b *Bot, m *discordgo.MessageCreate, msg string) {
			c, err := b.sess.State.Channel(m.ChannelID)
			if err != nil {
				logrus.Error(err)
				b.sess.ChannelMessageSend(m.ChannelID, "oops something went wrong, try again later")
				return
			}
			b.voiceMu.Lock()
			vh, ok := b.voiceHandlers[c.GuildID]
			if !ok {
				b.voiceMu.Unlock()
				b.sess.ChannelMessageSend(m.ChannelID, "Bot has to playing to repeat")
				return
			}
			vh.repeatChan <- 0
			b.voiceMu.Unlock()
		},
	},
}
