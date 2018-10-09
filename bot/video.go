package bot

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/fvdveen/mu2/services/search"
)

// Video represents a youtube video
type Video interface {
	OpusReader
	Name() string
	Author() string

	Announce() error

	ResetPlayback()
}

type video struct {
	v   *search.Video
	or  OpusReader
	ctx Context
	mu  sync.RWMutex

	opus []*frame
	done bool
	i    int
	len  int
}

type frame struct {
	o   []byte
	err error
}

// NewVideo creates a new video
func NewVideo(v *search.Video, or OpusReader, ctx Context) Video {
	vid := &video{
		v:   v,
		or:  or,
		ctx: ctx,
	}

	go vid.stream()

	return vid
}

func (v *video) OpusFrame() ([]byte, error) {
	var f *frame

	for {
		v.mu.RLock()
		if v.i < len(v.opus) {
			v.mu.RUnlock()
			break
		}
		v.mu.RUnlock()
	}

	v.mu.RLock()
	f = v.opus[v.i]
	v.i++
	if v.done && v.i >= len(v.opus) {
		v.i = 0
	}
	v.mu.RUnlock()

	return f.o, f.err
}

func (v *video) Author() string {
	return v.v.Author
}

func (v *video) Name() string {
	return v.v.Name
}

func (v *video) ResetPlayback() {
	v.i = 0
}

func (v *video) Announce() error {
	e := &discordgo.MessageEmbed{
		Description: fmt.Sprintf("**Now playing** [%s](%s)\nBy: %s", v.v.Name, v.v.URL, v.v.Author),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: v.v.ThumbnailURL,
		},
	}

	return v.ctx.SendEmbed(e)
}

func (v *video) stream() {
	for {
		o, err := v.or.OpusFrame()
		v.mu.Lock()
		v.opus = append(v.opus, &frame{
			o:   o,
			err: err,
		})
		v.len++
		v.mu.Unlock()
		if err != nil {
			v.done = true
			return
		}
	}
}
