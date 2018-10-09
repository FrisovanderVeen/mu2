package bot

import (
	"fmt"
	"io"

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

	done bool
}

// NewVideo creates a new video
func NewVideo(v *search.Video, or OpusReader, ctx Context) Video {
	return &video{
		v:   v,
		or:  or,
		ctx: ctx,
	}
}

func (v *video) OpusFrame() ([]byte, error) {
	if v.done {
		return nil, io.EOF
	}

	o, err := v.or.OpusFrame()
	if err == io.EOF {
		v.done = true
	}

	return o, err
}

func (v *video) Author() string {
	return v.v.Author
}

func (v *video) Name() string {
	return v.v.Name
}

func (v *video) ResetPlayback() {
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
