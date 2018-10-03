package events

import (
	"reflect"
	"strings"

	"github.com/fvdveen/mu2/config"
	"github.com/sirupsen/logrus"
)

const (
	Change EventType = iota
	Add
	Remove
)

type EventType uint8

type Event struct {
	// EventType shows what happened
	EventType EventType
	// The config key that got changes e.g. discord.token
	Key string

	Change    string
	Additions []string
	Removals  []string
	Database  config.Database
	Log       config.Log
}

// Watch puts all changes between the configs given by in into Events
func Watch(in <-chan *config.Config) <-chan *Event {
	ch := make(chan *Event)

	go func(in <-chan *config.Config, ch chan<- *Event) {
		last := &config.Config{}
		for conf := range in {
			if !reflect.DeepEqual(conf.Bot, last.Bot) {
				botChanges(ch, conf, last)
			}
			if !reflect.DeepEqual(conf.Log, last.Log) {
				logChanges(ch, conf, last)
			}
			if !reflect.DeepEqual(conf.Database, last.Database) {
				ch <- &Event{
					EventType: Change,
					Key:       "database",
					Change:    conf.Database.Type,
					Database:  conf.Database,
				}
			}

			last = conf
		}

		close(ch)
	}(in, ch)

	return ch
}

func logChanges(ch chan<- *Event, conf *config.Config, last *config.Config) {
	if conf.Log.Level != last.Log.Level {
		ch <- &Event{
			EventType: Change,
			Key:       "log.level",
			Change:    conf.Log.Level,
		}
	}
	if !reflect.DeepEqual(conf.Log.Discord, last.Log.Discord) {
		ch <- &Event{
			EventType: Change,
			Key:       "log.discord",
			Change:    "hook",
			Log:       conf.Log,
		}
	}
}

func botChanges(ch chan<- *Event, conf *config.Config, last *config.Config) {
	if conf.Bot.Discord.Token != last.Bot.Discord.Token {
		ch <- &Event{
			EventType: Change,
			Key:       "bot.discord.token",
			Change:    conf.Bot.Discord.Token,
		}
	}
	if conf.Bot.Prefix != last.Bot.Prefix {
		ch <- &Event{
			EventType: Change,
			Key:       "bot.prefix",
			Change:    conf.Bot.Prefix,
		}
	}

	if !reflect.DeepEqual(conf.Bot.Commands, last.Bot.Commands) {
		a, r := changes(conf.Bot.Commands, last.Bot.Commands)
		if len(a) == 0 && len(r) == 0 {
		} else if len(a) == 0 {
			ch <- &Event{
				EventType: Add,
				Key:       "bot.commands",
				Additions: a,
			}
		} else if len(r) == 0 {
			ch <- &Event{
				EventType: Remove,
				Key:       "bot.commands",
				Removals:  r,
			}
		} else {
			ch <- &Event{
				EventType: Change,
				Key:       "bot.commands",
				Additions: a,
				Removals:  r,
			}
		}
	}
}

func changes(new []string, old []string) (additions []string, removals []string) {
	oldComs := map[string]bool{}
	for _, com := range old {
		oldComs[com] = true
	}
	newComs := map[string]bool{}
	for _, com := range new {
		newComs[com] = true
	}

	for x := range oldComs {
		found := false
		for y := range newComs {
			if x == y {
				found = true
				break
			}
		}
		if found {
			continue
		}
		removals = append(removals, x)
	}

	for x := range newComs {
		double := false
		for y := range oldComs {
			if x == y {
				double = true
				break
			}
		}
		if double {
			continue
		}
		additions = append(additions, x)
	}

	return additions, removals
}

// Split splits the incoming event stream into bot, log and database based on the key of the event
func Split(in <-chan *Event) (<-chan *Event, <-chan *Event, <-chan *Event) {
	bot := make(chan *Event, 100)
	log := make(chan *Event, 100)
	db := make(chan *Event, 100)
	go func(in <-chan *Event, bot chan<- *Event, log chan<- *Event, db chan<- *Event) {
		for evnt := range in {
			d := strings.Split(evnt.Key, ".")[0]
			switch d {
			case "bot":
				bot <- evnt
			case "log":
				log <- evnt
			case "database":
				db <- evnt
			default:
				logrus.Warnf("Unknown event key: %s\n", d)
			}
		}

		close(bot)
		close(log)
		close(db)
	}(in, bot, log, db)

	return bot, log, db
}
