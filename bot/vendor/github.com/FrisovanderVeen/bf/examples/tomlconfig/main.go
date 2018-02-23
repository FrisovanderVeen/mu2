package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	bf "github.com/FrisovanderVeen/bf"
)

type tomlConfig struct {
	Discord discordConfig
}

type discordConfig struct {
	Prefix string
	Token  string
}

func main() {
	bot, err := bf.NewBot(decodeConfig("config.toml"))
	if err != nil {
		log.Fatalf("Could not create bot: %v", err)
	}

	if err := bot.Open(); err != nil {
		log.Fatalf("Could not open session: %v", err)
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	if err := bot.Close(); err != nil {
		log.Fatalf("Could not close session: %v", err)
	}
}

func decodeConfig(loc string) bf.OptionFunc {
	var conf tomlConfig
	if _, err := toml.DecodeFile(loc, &conf); err != nil {
		log.Printf("Could not decode config file: %v\n", err)
		return bf.EmptyOptionFunc
	}

	return func(b *bf.Bot) error {
		b.Token = conf.Discord.Token
		b.Prefix = conf.Discord.Prefix
		return nil
	}
}
