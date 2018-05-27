package bot

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/fvdveen/mu2/config"
)

// Bot represents a discord bot
type Bot struct {
	conf *config.Bot
	sess *discordgo.Session

	commands      []*command
	voiceHandlers map[string]*voiceHandler
	voiceMu       sync.RWMutex
}

// New creates a new bot
func New(conf *config.Bot) (*Bot, error) {
	b := &Bot{
		conf:          conf,
		voiceHandlers: make(map[string]*voiceHandler),
	}
	sess, err := discordgo.New("Bot " + conf.Token)
	if err != nil {
		return nil, err
	}
	b.sess = sess
	b.commands = commands

	b.initHandlers()

	return b, nil
}

// Open opens the discord session
func (b *Bot) Open() error {
	return b.sess.Open()
}

// Close closes the discord session
func (b *Bot) Close() error {
	for _, vh := range b.voiceHandlers {
		vh.stopChan <- 0
	}
	for len(b.voiceHandlers) != 0 {
	}
	return b.sess.Close()
}
