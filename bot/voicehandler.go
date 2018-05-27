package bot

import (
	"sync"
	"sync/atomic"

	"github.com/bwmarrin/discordgo"
	"github.com/fvdveen/mu2/pkg/queue"
	"github.com/sirupsen/logrus"
)

type voiceHandler struct {
	queue   *queue.Queue
	bot     *Bot
	sess    *discordgo.Session
	vc      *discordgo.VoiceConnection
	guildID string

	paused atomic.Value
	loop   atomic.Value
	repeat atomic.Value

	skipChan   chan interface{}
	stopChan   chan interface{}
	playChan   chan interface{}
	pauseChan  chan interface{}
	loopChan   chan interface{}
	repeatChan chan interface{}

	mu sync.RWMutex
	wg sync.WaitGroup
}

func newVoiceHandler(s *discordgo.Session, b *Bot, guildID string) *voiceHandler {
	v := &voiceHandler{
		queue:      queue.New(),
		bot:        b,
		guildID:    guildID,
		sess:       s,
		skipChan:   make(chan interface{}, 1),
		stopChan:   make(chan interface{}, 1),
		playChan:   make(chan interface{}, 1),
		pauseChan:  make(chan interface{}, 1),
		loopChan:   make(chan interface{}, 1),
		repeatChan: make(chan interface{}, 1),
	}
	v.paused.Store(false)
	v.loop.Store(false)
	v.repeat.Store(false)
	v.wg.Add(1)
	return v
}

func (vh *voiceHandler) handle(textChanID, voiceChanID, guildID string) {
	vh.wg.Wait()

	err := vh.joinVoiceChannel(guildID, voiceChanID)
	if err != nil {
		logrus.Errorf("[voiceHandler-handle] %v", err)
		vh.sess.ChannelMessageSend(textChanID, "something went wrong setting the speaking status")
		vh.bot.voiceMu.Lock()
		defer vh.bot.voiceMu.Unlock()
		err = vh.leaveVoiceChannel()
		if err != nil {
			logrus.Errorf("[voiceHandler-handle] %v", err)
		}
		return
	}

	for {
		vh.mu.Lock()
		if vh.queue == nil || vh.queue.Length() == 0 {
			vh.bot.voiceMu.Lock()
			defer vh.bot.voiceMu.Unlock()
			defer vh.mu.Unlock()

			err := vh.leaveVoiceChannel()
			if err != nil {
				logrus.Errorf("[voiceHandler-handle] %v", err)
			}
			return
		}
		vi, ok := vh.queue.PopFront().(*voiceItem)
		vh.mu.Unlock()
		if !ok {
			logrus.Error("[voiceHandler-handle] could not convert to voiceitem")
			continue
		}

		if vi.showMessage {
			vh.sess.ChannelMessageSend(vi.messageChannel, vi.message)
		}

		vh.playItem(vi)
		if vh.repeat.Load().(bool) {
			vh.addFront(vi)
		} else if vh.loop.Load().(bool) {
			vh.add(vi)
		}
	}
}

func (vh *voiceHandler) add(vi *voiceItem) {
	if vh.queue == nil {
		return
	}
	vh.mu.Lock()
	defer vh.mu.Unlock()
	vh.queue.PushBack(vi)
}

func (vh *voiceHandler) addFront(vi *voiceItem) {
	if vh.queue == nil {
		return
	}
	vh.mu.Lock()
	defer vh.mu.Unlock()
	vh.queue.PushFront(vi)
}

func (vh *voiceHandler) joinVoiceChannel(guildID, voiceChanID string) error {
	vc, err := vh.sess.ChannelVoiceJoin(guildID, voiceChanID, false, true)
	if err != nil {
		return err
	}

	err = vc.Speaking(true)
	if err != nil {
		return err
	}

	vh.vc = vc
	return nil
}

func (vh *voiceHandler) leaveVoiceChannel() error {
	delete(vh.bot.voiceHandlers, vh.guildID)
	err := vh.vc.Speaking(false)
	if err != nil {
		return err
	}

	return vh.vc.Disconnect()
}

func (vh *voiceHandler) playItem(vi *voiceItem) {
playloop:
	for _, f := range vi.data {
		for vh.paused.Load().(bool) {
			select {
			case <-vh.pauseChan:
				vh.paused.Store(true)
			case <-vh.playChan:
				vh.paused.Store(false)
			case <-vh.loopChan:
				vh.loop.Store(!vh.loop.Load().(bool))
			case <-vh.repeatChan:
				vh.repeat.Store(!vh.repeat.Load().(bool))
			case <-vh.skipChan:
				vh.paused.Store(false)
				break playloop
			case <-vh.stopChan:
				vh.mu.Lock()
				vh.queue = nil
				vh.mu.Unlock()
				vh.paused.Store(false)
				vh.loop.Store(false)
				vh.repeat.Store(false)
				return
			}
		}
		select {
		case <-vh.playChan:
			vh.paused.Store(false)
		case <-vh.pauseChan:
			vh.paused.Store(true)
		case <-vh.loopChan:
			vh.loop.Store(!vh.loop.Load().(bool))
		case <-vh.repeatChan:
			vh.repeat.Store(!vh.repeat.Load().(bool))
		case <-vh.skipChan:
			vh.paused.Store(false)
			break playloop
		case <-vh.stopChan:
			vh.mu.Lock()
			vh.queue = nil
			vh.mu.Unlock()
			vh.paused.Store(false)
			vh.loop.Store(false)
			vh.repeat.Store(false)
			return
		case vh.vc.OpusSend <- f:
		}
	}
}
