package queue

import (
	"fmt"
	"strconv"

	"github.com/fvdveen/mu2/bot"
	"github.com/fvdveen/mu2/commands"
	"github.com/sirupsen/logrus"
)

var _ = commands.Register(bot.NewCommand("queue", "queue displays the current queue", func(ctx bot.Context, args []string) {
	g, err := ctx.Guild()
	if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "queue"}).Errorf("Get guild: %v", err)
		return
	}

	vh, err := ctx.Bot().VoiceHandler(g.ID, "")
	if err != nil && err == bot.ErrVoiceStateNotFound {
		if err := ctx.Send("bot has to be playing to use queue"); err != nil {
			logrus.WithFields(map[string]interface{}{"type": "command", "command": "queue"}).Errorf("Send message: %v", err)
		}
		return
	} else if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "queue"}).Errorf("Get voice handler: %v", err)
		return
	}

	q := vh.Queue()

	var msg string

	for i, v := range q {
		msg = fmt.Sprintf("%s\n%d. %s - %s", msg, i+1, v.Name(), v.Author())
	}

	if msg == "" {
		msg = fmt.Sprintf("Queue is empty use `%splay` to add songs", ctx.Bot().Prefix())
	} else {
		msg = fmt.Sprintf("`%s`", msg)
	}

	if err := ctx.Send(msg); err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "queue"}).Errorf("Send embed: %v", err)
		return
	}
}))

var _ = commands.Register(bot.NewCommand("remove", "removes the queue item at the given position", func(ctx bot.Context, args []string) {
	g, err := ctx.Guild()
	if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "remove"}).Errorf("Get guild: %v", err)
		return
	}

	vh, err := ctx.Bot().VoiceHandler(g.ID, "")
	if err != nil && err == bot.ErrVoiceStateNotFound {
		if err := ctx.Send("bot has to be playing to use remove"); err != nil {
			logrus.WithFields(map[string]interface{}{"type": "command", "command": "remove"}).Errorf("Send message: %v", err)
		}
		return
	} else if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "remove"}).Errorf("Get voice handler: %v", err)
		return
	}

	for _, arg := range args {
		i, err := strconv.Atoi(arg)
		if err != nil {
			logrus.WithFields(map[string]interface{}{"type": "command", "command": "remove"}).Errorf("Convert string to int: %v", err)

			e, ok := err.(*strconv.NumError)
			if ok {
				if err := ctx.Send(fmt.Sprintf("Cannot use %s as integer", e.Num)); err != nil {
					logrus.WithFields(map[string]interface{}{"type": "command", "command": "remove"}).Errorf("Send message: %v", err)
				}
			}
			return
		}

		if err := vh.Remove(i - 1); err != nil {
			logrus.WithFields(map[string]interface{}{"type": "command", "command": "remove"}).Errorf("Remove item: %v", err)
		}
	}
}))

var _ = commands.Register(bot.NewCommand("reorder", "reorder", func(ctx bot.Context, args []string) {
	if len(args) < 2 {
		if err := ctx.Send("need at least 2 integers to swap position"); err != nil {
			logrus.WithFields(map[string]interface{}{"type": "command", "command": "reorder"}).Errorf("Send message: %v", err)
		}
	}

	i, err := strconv.Atoi(args[0])
	if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "reorder"}).Errorf("Convert string to int: %v", err)

		e, ok := err.(*strconv.NumError)
		if ok {
			if err := ctx.Send(fmt.Sprintf("Cannot use %s as integer", e.Num)); err != nil {
				logrus.WithFields(map[string]interface{}{"type": "command", "command": "reorder"}).Errorf("Send message: %v", err)
			}
		}
		return
	}

	j, err := strconv.Atoi(args[1])
	if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "reorder"}).Errorf("Convert string to int: %v", err)

		e, ok := err.(*strconv.NumError)
		if ok {
			if err := ctx.Send(fmt.Sprintf("Cannot use %s as integer", e.Num)); err != nil {
				logrus.WithFields(map[string]interface{}{"type": "command", "command": "reorder"}).Errorf("Send message: %v", err)
			}
		}
		return
	}

	i, j = i-1, j-1

	g, err := ctx.Guild()
	if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "reorder"}).Errorf("Get guild: %v", err)
		return
	}

	vh, err := ctx.Bot().VoiceHandler(g.ID, "")
	if err != nil && err == bot.ErrVoiceStateNotFound {
		if err := ctx.Send("bot has to be playing to use reorder"); err != nil {
			logrus.WithFields(map[string]interface{}{"type": "command", "command": "reorder"}).Errorf("Send message: %v", err)
		}
		return
	} else if err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "reorder"}).Errorf("Get voice handler: %v", err)
		return
	}

	if err := vh.Reorder(i, j); err != nil {
		logrus.WithFields(map[string]interface{}{"type": "command", "command": "reorder"}).Errorf("Reorder item: %v", err)
	}
}))
