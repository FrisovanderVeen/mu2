package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fvdveen/mu2/voice/youtube"

	"github.com/fvdveen/mu2/common"

	"github.com/fvdveen/mu2/bot"

	_ "github.com/fvdveen/mu2/commands/help"
	_ "github.com/fvdveen/mu2/commands/play"
)

func main() {
	f, err := os.Open("config.json")
	if err != nil {
		common.GetLogger().Error(err)
		return
	}

	defer func() {
		if err := f.Close(); err != nil {
			common.GetLogger().Error(err)
			return
		}
	}()

	if err := common.ReadConfig(f); err != nil {
		common.GetLogger().Error(err)
		return
	}
	common.SetupLogger(common.GetConfig())

	if err := youtube.Setup(); err != nil {
		common.GetLogger().Error(err)
		return
	}

	b, err := bot.New(
		bot.WithConfig(common.GetConfig()),
	)
	if err != nil {
		common.GetLogger().Error(err)
		return
	}

	b.Open()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	common.GetLogger().Info("Bot is now running press CRTL-C to exit")

	<-sc

	b.Close()
}
