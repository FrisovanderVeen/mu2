package bot

import (
	"encoding/binary"
	"io"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	buffer = make([][]byte, 0)
)

type voiceItem struct {
	data [][]byte

	messageChannel string
	message        string
	showMessage    bool
	voiceState     *discordgo.VoiceState
}

func (b *Bot) playSound(guildID string, vi *voiceItem) {
	b.voiceMu.Lock()
	defer b.voiceMu.Unlock()
	vh, ok := b.voiceHandlers[guildID]
	if !ok {
		b.voiceHandlers[guildID] = newVoiceHandler(b.sess, b, guildID)
		go b.voiceHandlers[guildID].handle(vi.messageChannel, vi.voiceState.ChannelID, guildID)
		b.voiceHandlers[guildID].add(vi)
		b.voiceHandlers[guildID].wg.Done()
		return
	}

	vh.add(vi)
}

func (b *Bot) playAirhorn(authorID, textChanID, guildID string) error {
	g, err := b.sess.State.Guild(guildID)
	if err != nil {
		return err
	}
	for _, vs := range g.VoiceStates {
		if vs.UserID == authorID {
			b.playSound(guildID, &voiceItem{
				data:           buffer,
				voiceState:     vs,
				messageChannel: textChanID,
				message:        "now playing some fucking airhorn or some shit",
				showMessage:    true,
			})
			return nil
		}
	}
	return ErrUnknownVoiceState
}

func loadSound() error {
	file, err := os.Open("airhorn.dca")
	if err != nil {
		return err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err = file.Close()
			return err
		}

		if err != nil {
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			return err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}
