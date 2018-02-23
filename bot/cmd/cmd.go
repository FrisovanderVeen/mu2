package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	bf "github.com/FrisovanderVeen/bf"
	"github.com/FrisovanderVeen/mu2/bot/commands"
	"github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/urfave/cli"
)

type tomlConfig struct {
	Discord discordConfig
}

type discordConfig struct {
	Prefix string
	Token  string
}

var globalFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "config, conf",
		Value: "config.toml",
		Usage: "The TOML settings file location",
	},
}

// NewApp returns an app which runs the bot and its endpoint
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Muze2"
	app.Usage = "A discord bot"
	app.Flags = globalFlags
	app.Action = func(c *cli.Context) error {
		loc := c.String("config")
		logger := logrus.New()

		bot, err := bf.NewBot(bf.ErrWriter(logger.WriterLevel(logrus.ErrorLevel)), decodeConfig(loc), bf.ErrPrefix(func() string { return time.Now().Format("15:04:05") }))
		if err != nil {
			return err
		}

		bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
			s.UpdateStatus(0, fmt.Sprintf("%shelp", bot.Prefix))
		})

		if err := bot.AddCommand(commands.Commands...); err != nil {
			return err
		}

		if err := bot.Open(); err != nil {
			return err
		}

		logger.Info("Logged in as")
		logger.Info(bot.Session.State.User.Username)
		logger.Info(bot.Session.State.User.ID)
		logger.Info("Bot is now running.  Press CTRL-C to exit.")
		logger.Info("-------------------")
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc

		if err := bot.Close(); err != nil {
			return err
		}

		return nil
	}

	return app
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
