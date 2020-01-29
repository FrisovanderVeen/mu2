package common

import (
	"os"

	"github.com/bwmarrin/discordgo"

	"github.com/fvdveen/mu2/config"
	"github.com/sirupsen/logrus"
)

var logger = &logrus.Logger{
	Out:          os.Stderr,
	Formatter:    &logrus.TextFormatter{},
	Hooks:        make(logrus.LevelHooks),
	Level:        logrus.InfoLevel,
	ExitFunc:     os.Exit,
	ReportCaller: true,
}

// Logger is the logger used in the bot
type Logger interface {
	WithField(string, interface{}) Logger
	WithFields(map[string]interface{}) Logger

	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Printf(string, ...interface{})
	Warnf(string, ...interface{})
	Warningf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	Panicf(string, ...interface{})

	Debug(...interface{})
	Info(...interface{})
	Print(...interface{})
	Warn(...interface{})
	Warning(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Panic(...interface{})
}

type logrusLogger struct {
	logrus.FieldLogger
}

func (l logrusLogger) WithField(s string, i interface{}) Logger {
	return logrusLogger{l.FieldLogger.WithField(s, i)}
}

func (l logrusLogger) WithFields(m map[string]interface{}) Logger {
	return logrusLogger{l.FieldLogger.WithFields(m)}
}

// GetLogger returns a preconfigured logger
func GetLogger() Logger {
	return logrusLogger{logger}
}

// GetCommandLogger returns the logger for command com
func GetCommandLogger(com string) Logger {
	return GetLogger().WithFields(map[string]interface{}{"command": com, "type": "command"})
}

// GetHandlerLogger returns the logger for handler hand
func GetHandlerLogger(hand string) Logger {
	return GetLogger().WithFields(map[string]interface{}{"handler": hand, "type": "handler"})
}

func GetVoiceHandlerLogger(conn *discordgo.VoiceConnection) Logger {
	return GetLogger().WithFields(map[string]interface{}{"guild-id": conn.GuildID, "channel-id": conn.ChannelID, "type": "voice-handler"})
}

// SetupLogger configures the logger
func SetupLogger(conf config.Config) {
	var lvl logrus.Level
	switch conf.Logger.Level {
	case "trace", "trc", "t":
		lvl = logrus.TraceLevel
	case "debug", "dbg", "d":
		lvl = logrus.DebugLevel
	case "info", "inf", "i":
		lvl = logrus.InfoLevel
	case "warning", "warn", "w":
		lvl = logrus.WarnLevel
	case "error", "err", "e":
		lvl = logrus.ErrorLevel
	case "fatal", "ftl", "f":
		lvl = logrus.FatalLevel
	case "panic", "pnc", "p":
		lvl = logrus.PanicLevel
	default:
		lvl = logrus.InfoLevel
		defer func() {
			GetLogger().Warnf("Unknown log level: %s", GetConfig().Logger.Level)
		}()
	}

	logger.SetLevel(lvl)

}
