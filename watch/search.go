package watch

import (
	"sync"
	
	"github.com/fvdveen/mu2/services/search/wrapper"
	"github.com/fvdveen/mu2/services/search"
	"github.com/fvdveen/mu2-config/events"
	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"github.com/micro/go-micro/client"
)

// SearchService watcher ch for search service events and applies them to the returned service
func SearchService(ch <-chan *events.Event, cc *api.Config, opts ...client.CallOption) (search.Service, <-chan interface{}) {
	d := make(chan interface{})
	s := wrapper.New(nil)

	var wg sync.WaitGroup

	wg.Add(1)

	go func(ch <-chan *events.Event, cc *api.Config, s wrapper.Service, d chan<- interface{}, opts ...client.CallOption) {
		logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "search-service"}).Debug("Starting...")
		var done = false
		for evnt := range ch {
			if evnt.Key != "services.search.location" {
				continue
			}

			srv := search.NewService(evnt.Change, cc, opts...)

			if err := s.SetService(srv); err != nil {
				logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "search-service"}).Errorf("Set service: %v", err)
				continue
			}

			logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "search-service"}).Debug("Set service")

			if !done {
				wg.Done()
				done = true
			}
		}

		logrus.WithFields(map[string]interface{}{"type": "watcher", "watcher": "search-service"}).Debug("Stopping...")
		close(d)
	}(ch, cc, s, d, opts...)

	wg.Wait()

	return s, d
}