package cmd

import (
	"github.com/fvdveen/mu2/config"
	"github.com/kz/discordrus"
	"github.com/sirupsen/logrus"
)

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
		logrus.Warnf("Unknown log level: %s using info instead", lvl)
		return logrus.InfoLevel
	}
}

func loadLogger(conf config.Log, log *logrus.Logger) {
	lvl := getLogLevel(conf.Level)
	if lvl == logrus.DebugLevel+1 {
		logrus.Warnf("Unknown log level: %v using info instead", lvl)
		lvl = logrus.InfoLevel
	}
	log.SetLevel(lvl)

	if conf.Discord.WebHook != "" {
		lvl = getLogLevel(conf.Discord.Level)
		log.AddHook(discordrus.NewHook(
			conf.Discord.WebHook,
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
