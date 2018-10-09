package voice

import (
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
}))

var _ = commands.Register(bot.NewCommand("resume", "resumes the currently playing track", func(ctx bot.Context, args []string) {
	vh, err := ctx.VoiceHandler()
	if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "resume"}).Errorf("Get voice handler: %v", err)
		return
	}
	vh.Resume()
}))

var _ = commands.Register(bot.NewCommand("skip", "skips the currently playing track", func(ctx bot.Context, args []string) {
	vh, err := ctx.VoiceHandler()
	if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "skip"}).Errorf("Get voice handler: %v", err)
		return
	}
	vh.Skip()
}))

var _ = commands.Register(bot.NewCommand("stop", "stops the playing of audio", func(ctx bot.Context, args []string) {
	vh, err := ctx.VoiceHandler()
	if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "stop"}).Errorf("Get voice handler: %v", err)
		return
	}
	vh.Stop()
}))

var _ = commands.Register(bot.NewCommand("loop", "loops the queue", func(ctx bot.Context, args []string) {
	vh, err := ctx.VoiceHandler()
	if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "loop"}).Errorf("Get voice handler: %v", err)
		return
	}
	vh.Loop()
}))

var _ = commands.Register(bot.NewCommand("repeat", "repeats the currently playing track", func(ctx bot.Context, args []string) {
	vh, err := ctx.VoiceHandler()
	if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "repeat"}).Errorf("Get voice handler: %v", err)
		return
	}
	vh.Repeat()
}))
