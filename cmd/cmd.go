package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	bf "github.com/FrisovanderVeen/bf"
	"github.com/FrisovanderVeen/mu2/commands"
	"github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/urfave/cli"
)

var globalFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "token, t",
		Value: "DGTOKEN",
		Usage: "The enviroment variable of the discord token",
	},
	cli.StringFlag{
		Name:  "prefix, p",
		Value: "DGPREFIX",
		Usage: "The enviroment variable of the discord prefix",
	},
}

// NewApp returns an app which runs the bot and its endpoint
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Muze2"
	app.Usage = "A discord bot"
	app.Flags = globalFlags
	app.Action = func(c *cli.Context) error {
		tokenEnv := c.String("token")
		prefixEnv := c.String("prefix")
		logger := logrus.New()

		bot, err := bf.NewBot(bf.ErrWriter(logger.WriterLevel(logrus.ErrorLevel)), getEnvVars(tokenEnv, prefixEnv), bf.ErrPrefix(func() string { return time.Now().Format("15:04:05") }))
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

func getEnvVars(tokenEnv string, prefixEnv string) bf.OptionFunc {
	token := os.Getenv(tokenEnv)
	prefix := os.Getenv(prefixEnv)
	return func(b *bf.Bot) error {
		b.Token = token
		b.Prefix = prefix

		return nil
	}
}
