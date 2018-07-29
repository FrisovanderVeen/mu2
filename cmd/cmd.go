package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fvdveen/mu2/db"

	"github.com/fvdveen/mu2/bot"
	"github.com/fvdveen/mu2/config"
	"github.com/joho/godotenv"
	"github.com/kz/discordrus"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"

	// Register dbs
	_ "github.com/fvdveen/mu2/db/postgres"
)

var globalFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:  "dotenv",
		Usage: "Load environment variables from a .env file",
	},
	&cli.StringFlag{
		Name:  "dotenv-loc",
		Usage: "Location of .env file (default: .env)",
	},
	&cli.StringFlag{
		Name:  "log-level",
		Usage: "Level of logging messages displayed",
		Value: "info",
	},
	&cli.BoolFlag{
		Name:  "dev",
		Usage: "Sets dotenv to true and the log-level to debug",
	},
}

func run(c *cli.Context) error {
	conf := config.Load()
	conf.Defaults()
	setupLogger(conf.Log)

	db, err := db.Get(conf.Database)
	if err != nil {
		logrus.Errorf("Could not connect to database: %v", err)
		return nil
	}

	b, err := bot.New(bot.WithConfig(*conf), bot.WithDB(db))
	if err != nil {
		logrus.Errorf("Could not create session: %v", err)
		return nil
	}

	if err := b.Run(); err != nil {
		logrus.Errorf("Could not run bot: %v", err)
		return nil
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	logrus.Info("Bot is now running press CRTL-C to exit")
	<-sc

	if err := b.Close(); err != nil {
		logrus.Errorf("Could not close bot: %v", err)
		return nil
	}

	return nil
}

func setup(c *cli.Context) error {
	if c.Bool("dotenv") || c.Bool("dev") {
		loc := ".env"
		if c.String("dotenv-loc") != "" {
			loc = c.String("dotenv-loc")
		}
		err := godotenv.Load(loc)
		if err != nil {
			logrus.Fatalf("Could not load .env file: %v", err)

			return nil
		}
	}

	lvl := c.String("log-level")
	if lvl != "" {
		err := os.Setenv("LOG_LEVEL", lvl)
		if err != nil {
			logrus.Fatalf("Could not set config value %s: %v", "LOG_LEVEL", err)

			return nil
		}
	}

	if c.Bool("dev") {
		err := os.Setenv("LOG_LEVEL", "debug")
		if err != nil {
			logrus.Fatalf("Could not set config value %s: %v", "LOG_LEVEL", err)

			return nil
		}
	}

	return nil
}

// New creates a new cli app
func New() *cli.App {
	app := &cli.App{
		Name:   "Mu2",
		Usage:  "A discord bot",
		Flags:  globalFlags,
		Before: setup,
		Action: run,
	}

	return app
}

func setupLogger(conf config.Log) {
	lvl := getLogLevel(conf.Level)
	if lvl == logrus.DebugLevel+1 {
		logrus.Warnf("Unknown log level: %v using info instead", lvl)
		lvl = logrus.InfoLevel
	}
	logrus.SetLevel(lvl)

	if conf.DiscordWebHook != "" {
		lvl = getLogLevel(conf.DiscordLevel)
		logrus.AddHook(discordrus.NewHook(
			conf.DiscordWebHook,
			lvl,
			&discordrus.Opts{
				EnableCustomColors: true,
				CustomLevelColors: &discordrus.LevelColors{
					Debug: 10170623,
					Info:  3581519,
					Warn:  14327864,
					Error: 13631488,
					Panic: 13631488,
					Fatal: 13631488,
				},
			},
		))
	}
}

func getLogLevel(lvl string) logrus.Level {
	switch lvl {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "err", "error":
		return logrus.ErrorLevel
	case "panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	default:
		return logrus.DebugLevel + 1
	}
}
