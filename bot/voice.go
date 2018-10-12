package bot

import (
	"fmt"
	"io"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const (
	eventResume voiceEventType = iota
	eventPause
	eventLoop
	eventRepeat
	eventSkip
	eventStop
)

type voiceEventType uint8

type voiceEvent struct {
	ret chan bool
	typ voiceEventType
}

// OpusReader returns an opus frame and an error
type OpusReader interface {
	OpusFrame() ([]byte, error)
}

// VoiceHandler is a wrapper around a voice connection
type VoiceHandler interface {
	Play(Video)
	Pause() bool
	Resume() bool
	Skip() bool
	Stop() bool
	Loop() bool
	Repeat() bool
	Queue() []Video
	Reorder(int, int) error
	Remove(int) error
}

type voiceHandler struct {
	c *discordgo.VoiceConnection
	b *bot

	mCID string

	q      *queue
	events chan *voiceEvent

	pause, loop, repeat, stop bool
}

func (b *bot) newVoiceHandler(gID string, cID string) (VoiceHandler, error) {
	if cID == "" {
		return nil, fmt.Errorf("can't joinc empty channel")
	}

	b.sessMu.RLock()
	c, err := b.sess.ChannelVoiceJoin(gID, cID, false, true)
	b.sessMu.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("join voice channel: %v", err)
	}

	v := &voiceHandler{
		c:      c,
		q:      newQueue(),
		b:      b,
		events: make(chan *voiceEvent),
	}

	go v.run()

	return v, nil
}

func (b *bot) VoiceHandler(gID string, cID string) (VoiceHandler, error) {
	var v VoiceHandler
	var ok bool

	b.voiceMu.Lock()
	defer b.voiceMu.Unlock()
	v, ok = b.voiceHandlers[gID]
	if !ok {
		var err error
		v, err = b.newVoiceHandler(gID, cID)
		if err != nil {
			return nil, err
		}

		b.voiceHandlers[gID] = v
	}

	return v, nil
}

func (vh *voiceHandler) Play(v Video) {
	vh.q.PushBack(v)
}

func (vh *voiceHandler) Queue() []Video {
	return vh.q.Copy()
}

func (vh *voiceHandler) Reorder(a int, b int) error {
	return vh.q.Reorder(a, b)
}

func (vh *voiceHandler) Remove(i int) error {
	return vh.q.Remove(i)
}

func (vh *voiceHandler) disconnect() {
	if err := vh.c.Disconnect(); err != nil {
		logrus.WithFields(map[string]interface{}{"type": "voice-handler", "guild": vh.c.GuildID}).Errorf("Close voice connection: %v", err)
	}

	vh.b.voiceMu.Lock()
	delete(vh.b.voiceHandlers, vh.c.GuildID)
	vh.b.voiceMu.Unlock()
}

func (vh *voiceHandler) run() {
	var v Video
	v = vh.first()

	for {
		if v == nil {
			vh.disconnect()
			return
		}

		if err := v.Announce(); err != nil {
			logrus.WithFields(map[string]interface{}{"type": "voice-handler", "guild": vh.c.GuildID}).Errorf("Announce video: %v", err)
		}

		if err := vh.playItem(v); err != nil {
			logrus.WithFields(map[string]interface{}{"type": "voice-handler", "guild": vh.c.GuildID}).Errorf("Play video: %v", err)
		}

		if vh.stop {
			vh.disconnect()
			return
		} else if vh.repeat {
			v.ResetPlayback()

			vh.q.PushFront(v)
		} else if vh.loop {
			v.ResetPlayback()

			vh.q.PushBack(v)
		}

		v = vh.q.PopFront()
	}
}

func (vh *voiceHandler) first() Video {
	for {
		v := vh.q.Front()
		if v != nil {
			break
		}
	}

	return vh.q.PopFront()
}

func (vh *voiceHandler) playItem(v Video) error {
	for {
		o, err := v.OpusFrame()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return fmt.Errorf("get opus: %v", err)
		}

		select {
		case evnt := <-vh.events:
			skip := vh.handleEvent(evnt)
			if skip {
				return nil
			}
		case vh.c.OpusSend <- o:
		}
	}
}

func (vh *voiceHandler) handleEvent(evnt *voiceEvent) bool {
	switch evnt.typ {
	case eventStop:
		vh.stop = true
		go func() {
			evnt.ret <- true
		}()
		return true
	case eventSkip:
		go func() {
			evnt.ret <- true
		}()
		return true
	case eventPause:
		go func() {
			evnt.ret <- true
		}()
		skip := vh.paused()
		if skip {
			return true
		}
	case eventResume:
		go func() {
			evnt.ret <- true
		}()
	case eventLoop:
		vh.loop = !vh.loop
		go func() {
			evnt.ret <- vh.loop
		}()
	case eventRepeat:
		vh.repeat = !vh.repeat
		go func() {
			evnt.ret <- vh.repeat
		}()
	}
	return false
}

func (vh *voiceHandler) paused() bool {
	for evnt := range vh.events {
		switch evnt.typ {
		case eventStop:
			vh.stop = true
			go func() {
				evnt.ret <- true
			}()
			return true
		case eventSkip:
			go func() {
				evnt.ret <- true
			}()
			return true
		case eventPause:
			go func() {
				evnt.ret <- true
			}()
		case eventResume:
			go func() {
				evnt.ret <- true
			}()
			return false
		case eventLoop:
			vh.loop = !vh.loop
			go func() {
				evnt.ret <- vh.loop
			}()
		case eventRepeat:
			vh.repeat = !vh.repeat
			go func() {
				evnt.ret <- vh.repeat
			}()
		}
	}
	return false
}

func (vh *voiceHandler) Pause() bool {
	ret := make(chan bool)
	select {
	case vh.events <- &voiceEvent{
		typ: eventPause,
		ret: ret,
	}:
	case <-time.After(time.Second):
	}
	return <-ret
}

func (vh *voiceHandler) Resume() bool {
	ret := make(chan bool)
	select {
	case vh.events <- &voiceEvent{
		typ: eventResume,
		ret: ret,
	}:
	case <-time.After(time.Second):
	}
	return <-ret
}

func (vh *voiceHandler) Skip() bool {
	ret := make(chan bool)
	select {
	case vh.events <- &voiceEvent{
		typ: eventSkip,
		ret: ret,
	}:
	case <-time.After(time.Second):
	}
	return <-ret
}

func (vh *voiceHandler) Stop() bool {
	ret := make(chan bool)
	select {
	case vh.events <- &voiceEvent{
		typ: eventStop,
		ret: ret,
	}:
	case <-time.After(time.Second):
	}
	return <-ret
}

func (vh *voiceHandler) Loop() bool {
	ret := make(chan bool)
	select {
	case vh.events <- &voiceEvent{
		typ: eventLoop,
		ret: ret,
	}:
	case <-time.After(time.Second):
	}
	return <-ret
}

func (vh *voiceHandler) Repeat() bool {
	ret := make(chan bool)
	select {
	case vh.events <- &voiceEvent{
		typ: eventRepeat,
		ret: ret,
	}:
	case <-time.After(time.Second):
	}
	return <-ret
}
