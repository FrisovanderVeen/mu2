package watch

import (
	"sync"

	"github.com/fvdveen/mu2/config/events"
	"github.com/fvdveen/mu2/log"
	"github.com/sirupsen/logrus"
)

// Log applies the changes from ch onto log
func Log(l *logrus.Logger, ch <-chan *events.Event, wg *sync.WaitGroup) <-chan interface{} {
	d := make(chan interface{})

	wg.Done()

	go func(l *logrus.Logger, ch <-chan *events.Event, d chan<- interface{}) {
		logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "log"}).Debug("Starting...")
		done := false
		var h log.Hook
		for !done {
			evnt, ok := <-ch
			if !ok {
				close(d)
				return
			}
			switch evnt.Key {
			case "log.level":
				l.SetLevel(log.GetLevel(evnt.Change))
				l.SetLevel(log.GetLevel(evnt.Change))
			case "log.discord":
				h = log.NewHookWrapper(log.DiscordHook(evnt.Log.Discord.Level, evnt.Log.Discord.WebHook))
				h.SetLevel(log.GetLevel(evnt.Log.Discord.Level))
				l.AddHook(h)
				logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "log"}).Debugf("Set discord hook and level to: %s, %s", evnt.Log.Discord.WebHook, evnt.Log.Discord.Level)
				done = true
			default:
				logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "log"}).Warnf("Unkonwn event key: %s", evnt.Key)
			}
		}
		for evnt := range ch {
			switch evnt.Key {
			case "log.level":
				l.SetLevel(log.GetLevel(evnt.Change))
				logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "log"}).Debugf("Set log level to: %s", evnt.Change)
			case "log.discord":
				h = log.NewHookWrapper(log.DiscordHook(evnt.Log.Discord.Level, evnt.Log.Discord.WebHook))
				logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "log"}).Debugf("Set discord hook and level to: %s, %s", evnt.Log.Discord.WebHook, evnt.Log.Discord.Level)
			default:
				logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "log"}).Warnf("Unknown event key: %s", evnt.Key)
			}
		}

		logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "log"}).Debug("Stopping...")

		close(d)
	}(l, ch, d)

	return d
}
