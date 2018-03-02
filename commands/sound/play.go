package sound

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fvdveen/bf"
	"github.com/fvdveen/mu2/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
	"google.golang.org/api/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
)

var (
	queries = make(map[string]*youtubeQuery)
	ytKey   = os.Getenv("YTKEY")
)

type youtubeVideo struct {
	ID     string
	Title  string
	Author string
}

type youtubeQuery struct {
	Videos []*youtubeVideo
}

var _ = commands.Register(bf.NewCommand(
	bf.Name("play"),
	bf.Trigger("play"),
	bf.Use("Queires youtube or plays the audio of the link"),
	bf.Action(func(ctx bf.Context) {
		if strings.HasPrefix(ctx.Message, "http://") || strings.HasPrefix(ctx.Message, "https://") {
			vidinf, err := ytdl.GetVideoInfo(ctx.Message)
			if err != nil {
				log.Errorf("could not get video info: %v", err)
				return
			}
			enqueue(vidinf, ctx)
			return
		}

		client := &http.Client{
			Transport: &transport.APIKey{Key: ytKey},
		}

		service, err := youtube.New(client)
		if err != nil {
			log.Errorf("Could not create new youtube client: %v", err)
			return
		}
		call := service.Search.List("id, snippet").
			Q(ctx.Message).
			MaxResults(5).
			Type("video")

		resp, err := call.Do()
		if err != nil {
			log.Errorf("Could not do youtube api call: %v", err)
			return
		}

		yq := &youtubeQuery{
			Videos: []*youtubeVideo{},
		}
		i := 1
		var embedItems []*discordgo.MessageEmbedField

		for _, video := range resp.Items {
			yv := &youtubeVideo{
				ID:     video.Id.VideoId,
				Title:  video.Snippet.Title,
				Author: video.Snippet.ChannelTitle,
			}
			yq.Videos = append(yq.Videos, yv)
			embedItems = append(embedItems, &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("%d. %s", i, yv.Title),
				Value: yv.Author,
			})
			i++
		}
		embed := &discordgo.MessageEmbed{
			Title:  "videos",
			Fields: embedItems,
		}
		if err := ctx.SendEmbed(embed); err != nil {
			log.Errorf("Could not send embed: %v", err)
			return
		}
		queries[ctx.MessageCreate.Author.ID] = yq
	}),
))

var _ = commands.Register(bf.NewCommand(
	bf.Name("choose"),
	bf.Trigger("choose"),
	bf.Use("Chooses a video from the query"),
	bf.Action(func(ctx bf.Context) {
		num, err := strconv.Atoi(ctx.Message)
		if err != nil {
			return
		}
		num--

		if num < 0 {
			return
		} else if num > len(queries[ctx.MessageCreate.Author.ID].Videos)-1 {
			return
		}

		vid := queries[ctx.MessageCreate.Author.ID].Videos[num]

		vidinf, err := ytdl.GetVideoInfoFromID(vid.ID)
		if err != nil {
			log.Errorf("could not get video info: %v", err)
			return
		}
		durl, err := vidinf.GetDownloadURL(vidinf.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0])
		if err != nil {
			log.Errorf("Could not get download url: %v", err)
			return
		}

		options := dca.StdEncodeOptions
		options.RawOutput = true
		options.Bitrate = 96
		options.Application = "lowdelay"
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

		if err := ctx.SendMessage(fmt.Sprintf("Added %s - %s to queue", vid.Title, vid.Author)); err != nil {
			log.Errorf("Could not send message: %v", err)
		}

		queue.Enqueue(sound{
			ctx:     ctx,
			view:    true,
			content: content,
			author:  vid.Author,
			name:    vid.Title,
		})
	}),
))

func enqueue(vidinf *ytdl.VideoInfo, ctx bf.Context) {
	durl, err := vidinf.GetDownloadURL(vidinf.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0])
	if err != nil {
		log.Errorf("Could not get download url: %v", err)
		return
	}

	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"
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

	if err := ctx.SendMessage(fmt.Sprintf("Added %s - %s to queue", vidinf.Title, vidinf.Author)); err != nil {
		log.Errorf("Could not send message: %v", err)
	}

	queue.Enqueue(sound{
		ctx:     ctx,
		author:  vidinf.Author,
		name:    vidinf.Title,
		view:    true,
		content: content,
	})
}
