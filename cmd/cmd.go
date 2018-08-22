package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kz/discordrus"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"

	"github.com/fvdveen/mu2/config"
	"github.com/fvdveen/mu2/db"

	// Register dbs
	_ "github.com/fvdveen/mu2/db/postgres"
)

var globalFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "token",
		Usage:   "discord token",
		EnvVars: []string{"DISCORD_TOKEN"},
	},
	&cli.StringFlag{
		Name:    "prefix",
		Usage:   "discord prefix",
		Value:   "$",
		EnvVars: []string{"DISCORD_PREFIX"},
	},
	&cli.StringFlag{
		Name:    "log-level",
		Usage:   "log level for stdout",
		EnvVars: []string{"LOG_LEVEL"},
	},
	&cli.StringFlag{
		Name:    "discord-webhook",
		Usage:   "discord webhook for logging",
		EnvVars: []string{"LOG_WEBHOOK_DISCORD"},
	},
	&cli.StringFlag{
		Name:    "discord-log-level",
		Usage:   "log level for discord",
		EnvVars: []string{"LOG_LEVEL_DISCORD"},
	},
	&cli.StringFlag{
		Name:    "db-host",
		Usage:   "host address for database",
		EnvVars: []string{"DB_HOST"},
	},
	&cli.StringFlag{
		Name:    "db-user",
		Usage:   "user for database",
		EnvVars: []string{"DB_USER"},
	},
	&cli.StringFlag{
		Name:    "db-password",
		Usage:   "password for database",
		EnvVars: []string{"DB_PASS"},
	},
	&cli.StringFlag{
		Name:    "db-ssl",
		Usage:   "ssl type for database",
		EnvVars: []string{"DB_SSL"},
	},
	&cli.StringFlag{
		Name:    "db-type",
		Usage:   "database type",
		EnvVars: []string{"DB_TYPE"},
	},
}

func run(c *cli.Context) error {
	conf := config.Load(c)

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

// New creates a new cli app
func New() *cli.App {
	app := &cli.App{
		Name:  "Mu2",
		Usage: "A discord bot",
		Flags: globalFlags,
		Before: func(c *cli.Context) error {
			return nil
		},
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
