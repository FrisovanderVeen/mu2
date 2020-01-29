package youtube

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"sync/atomic"
	"time"

	"github.com/fvdveen/mu2/common"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
)

var ErrVideoTooLong = errors.New("video too long")
var ErrNoSuchVideo = errors.New("no such video")

var maxTime int64 = math.MaxInt64

var youtubeService *youtube.Service

func MaxVideoLength() time.Duration {
	return time.Duration(atomic.LoadInt64(&maxTime))
}

func SetMaxVideoLength(d time.Duration) {
	atomic.StoreInt64(&maxTime, int64(d))
}

func Setup() error {
	key := common.GetConfig().Service.Youtube.APIKey
	srvc, err := youtube.NewService(context.Background(), option.WithAPIKey(key))
	if err != nil {
		return err
	}

	youtubeService = srvc
	return nil
}

type VideoInfo struct {
	URL   string
	Title string
}

func Search(title string) (*VideoInfo, error) {
	resp, err := youtubeService.Search.List("id,snippet").
		Context(context.Background()).
		Q(title).
		MaxResults(1).
		Type("video").
		Do()
	if err != nil {
		return nil, err
	}

	if len(resp.Items) == 0 {
		return nil, fmt.Errorf("%w: %s", ErrNoSuchVideo, title)
	}

	vi := &VideoInfo{
		URL:   fmt.Sprintf("https://www.youtube.com/watch?v=%s", resp.Items[0].Id.VideoId),
		Title: resp.Items[0].Snippet.Title,
	}
	return vi, nil
}

type VoiceItem struct {
	info    *VideoInfo
	encoder *Encoder
}

func NewVoiceItem(info *VideoInfo) (*VoiceItem, error) {
	vi := &VoiceItem{
		info: info,
	}

	d := NewDownloader(vi.info.URL)
	buf := &bytes.Buffer{}

	if err := d.Download(buf); err != nil {
		return nil, err
	}

	enc, err := NewEncoder(buf)
	if err != nil {
		return nil, err
	}

	vi.encoder = enc

	return vi, nil
}

func (vi *VoiceItem) Title() string {
	return vi.info.Title
}

func (vi *VoiceItem) OpusFrame() ([]byte, error) {
	return vi.encoder.OpusFrame()
}

type Downloader struct {
	url string
}

func NewDownloader(url string) *Downloader {
	d := &Downloader{
		url: url,
	}
	return d
}

func (d *Downloader) Download(dest io.Writer) error {
	videoInfo, err := ytdl.GetVideoInfo(d.url)
	if err != nil {
		return err
	}

	if videoInfo.Duration > MaxVideoLength() {
		return fmt.Errorf("%w: %v", ErrVideoTooLong, videoInfo.Duration)
	}

	format := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0]
	err = videoInfo.Download(format, dest)
	if err != nil {
		return err
	}

	return nil
}

type Encoder struct {
	encodeSession *dca.EncodeSession
}

func NewEncoder(r io.Reader) (*Encoder, error) {
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = dca.AudioApplicationLowDelay
	options.Channels = 2
	options.FrameRate = 48000
	options.FrameDuration = 20
	options.CompressionLevel = 1
	options.PacketLoss = 1
	options.BufferedFrames = 200
	options.VBR = true
	options.Threads = 4

	es, err := dca.EncodeMem(r, options)
	if err != nil {
		return nil, err
	}

	e := &Encoder{
		encodeSession: es,
	}

	return e, nil
}

func (e *Encoder) OpusFrame() ([]byte, error) {
	return e.encodeSession.OpusFrame()
}

func (e *Encoder) Close() error {
	e.encodeSession.Cleanup()
	return nil
}
