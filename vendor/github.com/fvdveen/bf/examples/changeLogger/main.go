package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/fvdveen/bf"
)

var test = &bf.NewCommand(
	bf.Name("test"),
	bf.Trigger("test"),
	bf.Use("Trows a error containing the message"),
	bf.Action(func(ctx bf.Context) {
		ctx.Error(errors.New(ctx.Message))
	}),
)

func main() {
	logger := logrus.New()
	bot, err := bf.NewBot(bf.Token("TOKEN"), bf.Prefix("-"), bf.ErrWriter(logger.Writer()), bf.ErrPrefix(func() string { return time.Now().String() }))
	if err != nil {
		logrus.Fatal(err)
	}
	bot.AddCommand(test)

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
