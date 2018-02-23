package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/FrisovanderVeen/bf"
)

var ping = &bf.Command{
	Name:    "ping",
	Trigger: "ping",
	Use:     "sends pong to the text channel",
	Action: func(ctx bf.Context) {
		if err := ctx.SendMessage("ping"); err != nil {
			ctx.Error(err)
		}
	},
}

var pong = &bf.Command{
	Name:    "pong",
	Trigger: "pong",
	Use:     "sends ping to the text channel",
	Action: func(ctx bf.Context) {
		if err := ctx.SendMessage("ping"); err != nil {
			ctx.Error(err)
		}
	},
}

func main() {
	bot, err := bf.NewBot(bf.Token("TOKEN"), bf.Prefix("-"))
	if err != nil {
		log.Printf("%v", err)
	}
	bot.AddCommand(ping, pong)

	if err := bot.Open(); err != nil {
		bot.Error(err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	if err := bot.Close(); err != nil {
		bot.Error(err)
		return
	}
}
