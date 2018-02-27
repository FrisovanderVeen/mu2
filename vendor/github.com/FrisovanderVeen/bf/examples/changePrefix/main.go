package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	bf "github.com/FrisovanderVeen/bf"
)

var ping = bf.NewCommand(
	bf.Name("ping"),
	bf.Trigger("ping"),
	bf.Use("Sends pong to the text channel"),
	bf.Action(func(ctx bf.Context) {
		if err := ctx.SendMessage("pong"); err != nil {
			ctx.Error(err)
		}
	}),
)

var pong = bf.NewCommand(
	bf.Name("pong"),
	bf.Trigger("pong"),
	bf.Use("Sends ping to the text channel"),
	bf.Action(func(ctx bf.Context) {
		if err := ctx.SendMessage("ping"); err != nil {
			ctx.Error(err)
		}
	}),
)

var cp = bf.NewCommand(
	bf.Name("change prefix"),
	bf.Trigger("prefix"),
	bf.Use("Chages the prefix for commands"),
	bf.Action(func(ctx bf.Context) {
		if strings.HasPrefix(ctx.Message, "prefix") {
			ctx.Message = strings.TrimPrefix(ctx.Message, "prefix")
		} else {
			return
		}
		if strings.HasPrefix(ctx.Message, " ") {
			ctx.Message = strings.TrimPrefix(ctx.Message, " ")
		}
		ctx.Bot.Prefix = ctx.Message
	}),
)

func main() {
	bot, err := bf.NewBot(bf.Token("TOKEN"), bf.Prefix("-"))
	if err != nil {
		log.Printf("%v", err)
	}
	bot.AddCommand(ping, pong, cp)

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
