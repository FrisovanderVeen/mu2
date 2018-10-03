package watch

import (
	"sync"

	"github.com/fvdveen/mu2/config/events"
	"github.com/fvdveen/mu2/db"
	"github.com/fvdveen/mu2/db/wrapper"
	"github.com/sirupsen/logrus"
)

// DB applies the changes from ch onto the service returned
func DB(ch <-chan *events.Event, wg *sync.WaitGroup) (wrapper.Service, <-chan interface{}) {
	done := false
	s := wrapper.New(nil)

	logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "db"}).Debug("Starting...")

	for !done {
		evnt := <-ch
		if evnt.Key != "database" {
			logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "db"}).Warnf("Unkown event key: %s", evnt.Key)
			continue
		}
		srv, err := db.Get(evnt.Database)
		if err != nil {
			logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "db"}).Errorf("Watch database: %v", err)
			continue
		}
		s.SetService(srv)
		logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "db"}).Debugf("Set database to: %s", evnt.Database.Type)
		done = true
	}

	wg.Done()

	d := make(chan interface{})
	go func(s wrapper.Service, ch <-chan *events.Event, d chan<- interface{}) {
		for evnt := range ch {
			if evnt.Key != "database" {
				logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "db"}).Warnf("Unkown event key: %s", evnt.Key)
				continue
			}
			srv, err := db.Get(evnt.Database)
			if err != nil {
				logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "db"}).Errorf("Watch database: %v", err)
				continue
			}
			s.SetService(srv)
			logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "db"}).Debugf("Set database to: %s", evnt.Database.Type)
		}

		logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "db"}).Debug("Stopping...")

		close(d)
	}(s, ch, d)

	return s, d
}
