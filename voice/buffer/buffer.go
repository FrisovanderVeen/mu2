package buffer

import (
	"io"
	"sync"

	"github.com/fvdveen/mu2/bot"
)

type opusFrame struct {
	frame []byte
	err   error
}

type Buffer struct {
	frames []opusFrame
	sub    interface {
		bot.OpusPlayer
		bot.VideoInfo
	}
	mu    sync.Mutex
	pos   int
	done  bool
	async bool
}

type OptionFunc func(*Buffer)

func WithAsync() func(*Buffer) {
	return func(b *Buffer) {
		b.async = true
	}
}

func New(enc interface {
	bot.OpusPlayer
	bot.VideoInfo
}, opts ...OptionFunc) *Buffer {
	b := &Buffer{
		frames: []opusFrame{},
		sub:    enc,
	}

	for _, opt := range opts {
		opt(b)
	}

	if b.async {
		b.readInAsync()
	}

	return b
}

func (b *Buffer) lazyOpusFrame() ([]byte, error) {
	if b.pos < len(b.frames) {
		frame := b.frames[b.pos]
		if frame.err == io.EOF {
			b.pos = 0
		} else {
			b.pos++
		}
		return frame.frame, frame.err
	}

	b.pos = 0
	return nil, io.EOF
}

func (b *Buffer) OpusFrame() ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.done {
		return b.lazyOpusFrame()
	}
	if b.pos < len(b.frames) {
		frame := b.frames[b.pos]
		b.pos++
		return frame.frame, frame.err
	}

	b.pos++
	bytes, err := b.sub.OpusFrame()
	if err == io.EOF {
		b.pos = 0
		b.done = true
		return nil, io.EOF
	}
	b.frames = append(b.frames, opusFrame{bytes, err})
	return bytes, err
}

func (b *Buffer) ResetPlayback() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.pos = 0
}

func (b *Buffer) Title() string {
	return b.sub.Title()
}

func (b *Buffer) Close() error {
	c, ok := b.sub.(io.Closer)
	if ok {
		return c.Close()
	}

	return nil
}

func (b *Buffer) readInAsync() {
	go func() {
		for {
			b.mu.Lock()
			if b.done {
				return
			}
			bytes, err := b.sub.OpusFrame()
			if err == io.EOF {
				b.done = true
				b.frames = append(b.frames, opusFrame{bytes, err})
				b.mu.Unlock()
				return
			}
			b.frames = append(b.frames, opusFrame{bytes, err})
			b.mu.Unlock()
		}
	}()
}
