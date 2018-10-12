package voice

import (
	"fmt"

	"github.com/fvdveen/mu2/bot"
	"github.com/fvdveen/mu2/commands"
	"github.com/sirupsen/logrus"
)

var _ = commands.Register(bot.NewCommand("pause", "pauses the currently playing track", func(ctx bot.Context, args []string) {
	vh, err := ctx.VoiceHandler()
	if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "pause"}).Errorf("Get voice handler: %v", err)
		return
	}
	vh.Pause()
	if err := ctx.Send(fmt.Sprintf("Paused playing")); err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "pause"}).Errorf("Send message: %v", err)
	}
}))

var _ = commands.Register(bot.NewCommand("resume", "resumes the currently playing track", func(ctx bot.Context, args []string) {
	vh, err := ctx.VoiceHandler()
	if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "resume"}).Errorf("Get voice handler: %v", err)
		return
	}
	vh.Resume()
	if err := ctx.Send(fmt.Sprintf("Resumed playing")); err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "resume"}).Errorf("Send message: %v", err)
	}
}))

var _ = commands.Register(bot.NewCommand("skip", "skips the currently playing track", func(ctx bot.Context, args []string) {
	vh, err := ctx.VoiceHandler()
	if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "skip"}).Errorf("Get voice handler: %v", err)
		return
	}
	vh.Skip()
	if err := ctx.Send(fmt.Sprintf("Skipped a song")); err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "skip"}).Errorf("Send message: %v", err)
	}
}))

var _ = commands.Register(bot.NewCommand("stop", "stops the playing of audio", func(ctx bot.Context, args []string) {
	vh, err := ctx.VoiceHandler()
	if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "stop"}).Errorf("Get voice handler: %v", err)
		return
	}
	vh.Stop()
	if err := ctx.Send(fmt.Sprintf("Stopped playing audio")); err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "stop"}).Errorf("Send message: %v", err)
	}
}))

var _ = commands.Register(bot.NewCommand("loop", "loops the queue", func(ctx bot.Context, args []string) {
	vh, err := ctx.VoiceHandler()
	if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "loop"}).Errorf("Get voice handler: %v", err)
		return
	}
	l := vh.Loop()
	if err := ctx.Send(fmt.Sprintf("Loop set to `%t`", l)); err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "loop"}).Errorf("Send message: %v", err)
	}
}))

var _ = commands.Register(bot.NewCommand("repeat", "repeats the currently playing track", func(ctx bot.Context, args []string) {
	vh, err := ctx.VoiceHandler()
	if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "repeat"}).Errorf("Get voice handler: %v", err)
		return
	}
	r := vh.Repeat()
	if err := ctx.Send(fmt.Sprintf("Repeat set to `%t`", r)); err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "repeat"}).Errorf("Send message: %v", err)
	}
}))
