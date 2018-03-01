package sound

import (
	"github.com/FrisovanderVeen/bf"
	"github.com/FrisovanderVeen/mu2/commands"
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
)

var _ = commands.Register(bf.NewCommand(
	bf.Name("play"),
	bf.Trigger("play"),
	bf.Use("Plays the audio of the link"),
	bf.Action(func(ctx bf.Context) {
		options := dca.StdEncodeOptions
		options.RawOutput = true
		options.Bitrate = 96
		options.Application = "lowdelay"
		vidinf, err := ytdl.GetVideoInfo(ctx.Message)
		if err != nil {
			log.Errorf("could not get video info: %v", err)
			return
		}

		durl, err := vidinf.GetDownloadURL(vidinf.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0])
		if err != nil {
			log.Errorf("Could not get download url: %v", err)
			return
		}

		encSess, err := dca.EncodeFile(durl.String(), options)
		if err != nil {
			log.Errorf("Could not encode %s: %v", durl.String(), err)
			return
		}
		defer encSess.Cleanup()

		content, err := encodeSessionToBytes(encSess)
		if err != nil {
			log.Errorf("Could not get bytes from encode session: %v", err)
			return
		}

		queue.Enqueue(sound{
			ctx:     ctx,
			view:    false,
			content: content,
		})
	}),
))
