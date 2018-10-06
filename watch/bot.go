package watch

import (
	"sync"

	"github.com/fvdveen/mu2-config"
	"github.com/fvdveen/mu2-config/events"
	"github.com/fvdveen/mu2/bot"
	"github.com/fvdveen/mu2/commands"
	"github.com/fvdveen/mu2/db"
	"github.com/sirupsen/logrus"
)

// Bot creates a bot and watches the channel for changes to the bot
// if anything is sent on ping a healthcheck will be done and the response will be sent on the error channel
func Bot(ch <-chan *events.Event, ping <-chan interface{}, s db.Service, wg *sync.WaitGroup) (<-chan interface{}, <-chan error) {
	var b bot.Bot
	d := make(chan interface{})

	go func(ch <-chan *events.Event, d chan<- interface{}, s db.Service) {
		logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "bot"}).Debug("Starting...")
		b = getBot(ch, s)
		var (
			token       string
			updateToken = false
		)
		var (
			prefix       string
			updatePrefix = false
		)

		wg.Done()

		for evnt := range ch {
			switch evnt.Key {
			case "bot.discord.token":
				logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Debugf("Changing token")
				token = evnt.Change
				if err := b.SetToken(token); err != nil {
					logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "bot"}).Errorf("Set token: %v", err)
					updateToken = true
					continue
				}
				updateToken = false
			case "bot.prefix":
				logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Debugf("Changing prefix")
				prefix = evnt.Change
				if err := b.SetPrefix(prefix); err != nil {
					logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "bot"}).Errorf("Set token: %v", err)
					updatePrefix = true
					continue
				}
				updatePrefix = false
			case "bot.commands":
				switch evnt.EventType {
				case events.Change:
					logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Debugf("Updating commands")
					for _, c := range evnt.Removals {
						if err := b.RemoveCommand(c); err != nil {
							logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Errorf("Remove command: %v", err)
							continue
						}
					}
					for _, cmd := range evnt.Additions {
						c, err := commands.Get(cmd)
						if err != nil {
							logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Warnf("%v", err)
							continue
						}
						if err := b.AddCommand(c); err != nil {
							logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Errorf("Add command: %v", err)
						}
					}
				case events.Remove:
					logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Debugf("Removing commands")
					for _, c := range evnt.Removals {
						if err := b.RemoveCommand(c); err != nil {
							logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Errorf("Remove command: %v", err)
							continue
						}
					}
				case events.Add:
					logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Debugf("Adding commands")
					for _, cmd := range evnt.Additions {
						c, err := commands.Get(cmd)
						if err != nil {
							logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Warnf("%v", err)
							continue
						}
						if err := b.AddCommand(c); err != nil {
							logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Errorf("Add command: %v", err)
						}
					}
				}
			}
			if updateToken {
				if err := b.SetToken(token); err != nil {
					logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "bot"}).Errorf("Set token: %v", err)
					updateToken = true
					continue
				}
				updateToken = false
			}
			if updatePrefix {
				if err := b.SetPrefix(prefix); err != nil {
					logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "bot"}).Errorf("Set token: %v", err)
					updatePrefix = true
					continue
				}
				updatePrefix = false
			}
		}
		logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "bot"}).Debug("Stopping...")

		if err := b.Close(); err != nil {
			logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Errorf("Close bot: %v", err)
		}

		close(d)
	}(ch, d, s)

	wg.Wait()

	res := make(chan error)

	go func(ping <-chan interface{}, res chan<- error) {
		for range ping {
			res <- b.Ping()
		}

		close(res)
	}(ping, res)

	return d, res
}

// helper function for Watch
func getBot(ch <-chan *events.Event, s db.Service) bot.Bot {
	var (
		b      bot.Bot
		err    error
		token  string
		prefix string
	)
	cmds := map[string]bool{}
	done := false

	for !done {
		evnt, ok := <-ch
		if !ok {
			return b
		}

		switch evnt.Key {
		case "bot.discord.token":
			logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Debugf("Changing token")
			token = evnt.Change
		case "bot.prefix":
			logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Debugf("Changing prefix")
			prefix = evnt.Change
		case "bot.commands":
			switch evnt.EventType {
			case events.Change:
				logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Debugf("Updating commands")
				for _, c := range evnt.Removals {
					delete(cmds, c)
				}
				for _, c := range evnt.Additions {
					cmds[c] = true
				}
			case events.Remove:
				logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Debugf("Removing commands")
				for _, c := range evnt.Removals {
					delete(cmds, c)
				}
			case events.Add:
				logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Debugf("Adding commands")
				for _, c := range evnt.Additions {
					cmds[c] = true
				}
			}
		}

		if token == "" {
			continue
		}

		logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Debug("Creating bot")

		b, err = bot.New(bot.WithConfig(config.Bot{
			Discord: config.Discord{
				Token: token,
			},
			Prefix: prefix,
		}), bot.WithDB(s))
		if err != nil {
			logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Errorf("Create bot: %v", err)
			continue
		}

		logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Debug("Opening session")

		if err := b.Open(); err != nil {
			logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Errorf("Open bot session: %v", err)
			continue
		}

		logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Debug("Adding commands")

		for cmd := range cmds {
			c, err := commands.Get(cmd)
			if err != nil {
				logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Warnf("%v", err)
				continue
			}

			if err := b.AddCommand(c); err != nil {
				logrus.WithFields(map[string]interface{}{"type": "watch", "watch": "bot"}).Errorf("Add command: %v", err)
			}
		}

		done = true
	}

	return b
}
