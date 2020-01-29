package bot

import (
	"errors"
	"io"
	"sync"

	"github.com/fvdveen/mu2/common"

	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/go-multierror"
)

var (
	// ErrVoiceStateNotFound is returned when the voice state was not found
	ErrVoiceStateNotFound = errors.New("voice state not found")
	errStopPlaying        = errors.New("voice handler recieved stop signal")
)

type OpusPlayer interface {
	OpusFrame() ([]byte, error)
}

type ResetOpusPlayer interface {
	OpusPlayer
	ResetPlayback()
}

type VideoInfo interface {
	Title() string
}

// VoiceItem is something that can be played over voice
type VoiceItem interface {
	ResetOpusPlayer
	VideoInfo
}

// VoiceHandler handles voice connections for the bot
type VoiceHandler struct {
	bot *Bot

	conn *discordgo.VoiceConnection

	start sync.Once

	queue *queue

	curPlaying   VoiceItem
	curPlayingMu sync.RWMutex

	repeat bool
}

// GetVoiceState returns the voice state of the user who triggered the command
func (ctx *Context) GetVoiceState() (*discordgo.VoiceState, error) {
	g, err := ctx.Session.Guild(ctx.MessageCreate.GuildID)
	if err != nil {
		return nil, err
	}

	for _, vs := range g.VoiceStates {
		if vs.UserID == ctx.MessageCreate.Author.ID {
			return vs, nil
		}
	}
	return nil, ErrVoiceStateNotFound
}

// GetVoiceHandler returns the voice handler for the guild of the user who triggered the command
func (ctx *Context) GetVoiceHandler() (*VoiceHandler, error) {
	ctx.Bot.voiceHandlersMu.RLock()
	vh, ok := ctx.Bot.voiceHandlers[ctx.MessageCreate.GuildID]
	ctx.Bot.voiceHandlersMu.RUnlock()
	if !ok {
		ctx.Bot.voiceHandlersMu.Lock()
		defer ctx.Bot.voiceHandlersMu.Unlock()

		vh, ok = ctx.Bot.voiceHandlers[ctx.MessageCreate.GuildID]
		if !ok {
			vs, err := ctx.GetVoiceState()
			if err != nil {
				return nil, err
			}
			vh, err = ctx.Bot.newVoiceHandler(vs.GuildID, vs.ChannelID)
			if err != nil {
				return nil, err
			}
		}
	}
	return vh, nil
}

// newVoiceHandler creates a new voice handler
// the caller has to lock voiceHandlerMu to prevent data races
func (b *Bot) newVoiceHandler(guildID, channelID string) (*VoiceHandler, error) {
	vc, err := b.session.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return nil, err
	}

	if err := vc.Speaking(true); err != nil {
		if erro := vc.Disconnect(); err != nil {
			return nil, multierror.Append(err, erro)
		}
		return nil, err
	}

	vh := &VoiceHandler{
		bot:   b,
		conn:  vc,
		queue: new(queue),
	}

	b.voiceHandlers[guildID] = vh

	return vh, nil
}

func (vh *VoiceHandler) CurrentPlaying() VoiceItem {
	vh.curPlayingMu.RLock()
	vh.curPlayingMu.RUnlock()
	return vh.curPlaying
}

// Play adds vi to the queue
func (vh *VoiceHandler) Play(vi VoiceItem) {
	vh.queue.Add(vi)
	vh.start.Do(func() {
		go vh.run()
	})
}

func (vh *VoiceHandler) run() {
	for {
		if !vh.repeat {
			vh.curPlayingMu.Lock()
			vh.curPlaying = vh.queue.Next()
			vh.curPlayingMu.Unlock()
		} else {
			vh.curPlayingMu.RLock()
			vh.curPlaying.ResetPlayback()
			vh.curPlayingMu.RUnlock()
		}

		vh.curPlayingMu.Lock()
		if vh.curPlaying == nil {
			if err := vh.Close(); err != nil {
				common.GetVoiceHandlerLogger(vh.conn).Errorf("close voice handler: %v", err)
				break
			}
			break
		}
		vh.curPlayingMu.Unlock()

		if err := vh.play(); err != nil && errors.Is(err, errStopPlaying) {
			break
		} else if err != nil {
			common.GetVoiceHandlerLogger(vh.conn).Errorf("play voice item: %v", err)
			break
		}
	}
}

func (vh *VoiceHandler) play() error {
	vh.curPlayingMu.RLock()
	defer vh.curPlayingMu.RUnlock()

	for {
		frame, err := vh.curPlaying.OpusFrame()
		if err != nil && errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return err
		}

		vh.conn.OpusSend <- frame
	}
}

func (vh *VoiceHandler) Close() error {
	if err := vh.conn.Speaking(false); err != nil {
		return err
	}
	vh.conn.Close()
	return nil
}
