package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	bf "github.com/FrisovanderVeen/bf"
	"github.com/bwmarrin/discordgo"
	"github.com/op/go-logging"
	"github.com/urfave/cli"

	"github.com/FrisovanderVeen/mu2/commands"
	_ "github.com/FrisovanderVeen/mu2/commands/help"
	_ "github.com/FrisovanderVeen/mu2/commands/info"
	_ "github.com/FrisovanderVeen/mu2/commands/pingpong"
)

var log = logging.MustGetLogger("cmd")

var globalFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "token, t",
		Value: "DGTOKEN",
		Usage: "The environment variable of the discord token",
	},
	cli.StringFlag{
		Name:  "prefix, p",
		Value: "DGPREFIX",
		Usage: "The environment variable of the discord prefix",
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
		logging.SetBackend(logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{module} ▶ %{level:.4s} %{id:03x} %{message} %{color:reset}`)))

		bot, err := bf.NewBot(getEnvVars(tokenEnv, prefixEnv))
		if err != nil {
			log.Critical("Could not make bot: %v", err)
			return err
		}

		bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
			s.UpdateStatus(0, fmt.Sprintf("%shelp", bot.Prefix))
		})

		if err := bot.AddCommand(commands.Commands...); err != nil {
			log.Errorf("Could not add command: %v", err)
			return err
		}

		if err := bot.Open(); err != nil {
			log.Critical("Could not open session: %v", err)
			return err
		}

		log.Info("Logged in as")
		log.Info(bot.Session.State.User.Username)
		log.Info(bot.Session.State.User.ID)
		log.Info("Bot is now running. Press CTRL-C to exit.")
		log.Info("-------------------")
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc

		if err := bot.Close(); err != nil {
			log.Errorf("Could not close session: %v", err)
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
