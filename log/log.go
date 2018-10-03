package log

import (
	"errors"
	"sync"

	"github.com/kz/discordrus"
	"github.com/sirupsen/logrus"
)

// Hook is a wrapper around a logrus hook with interchangable hook
type Hook interface {
	logrus.Hook
	// SetHook wil return a non-nil error if the hook is nil
	SetHook(logrus.Hook) error
	// SetLevel sets the minimum level for logs to be sent
	SetLevel(logrus.Level)
}

// GetLevel returns the logrus level from the given string
// It returns logrus.InfoLevel if the string was unrecognised
func GetLevel(lvl string) logrus.Level {
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

// DiscordHook creates a new logrus discord hook
func DiscordHook(lvl, url string) logrus.Hook {
	return discordrus.NewHook(
		url,
		GetLevel(lvl),
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
	)
}

// NewHookWrapper creates a new hook wrapper
func NewHookWrapper(h logrus.Hook) Hook {
	return &hook{
		h:   h,
		lvl: logrus.PanicLevel,
	}
}

type hook struct {
	h   logrus.Hook
	mu  sync.RWMutex
	lvl logrus.Level
}

func (h *hook) Fire(e *logrus.Entry) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if e.Level > h.lvl {
		return nil
	}
	return h.h.Fire(e)
}

func (h *hook) Levels() []logrus.Level {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.h.Levels()
}

func (h *hook) SetHook(hook logrus.Hook) error {
	if h == nil {
		return errors.New("given hook is nil")
	}

	h.mu.Lock()
	h.h = hook
	h.mu.Unlock()

	return nil
}

func (h *hook) SetLevel(lvl logrus.Level) {
	h.mu.Lock()
	h.lvl = lvl
	h.mu.Unlock()
}
