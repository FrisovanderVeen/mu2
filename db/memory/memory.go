package memory

import (
	"sync"

	"github.com/fvdveen/mu2/config"
	"github.com/fvdveen/mu2/db"
)

var _ = db.Register("memory", New)

// New creates a in-memory learnable db
func New(conf config.Database) (db.Service, error) {
	s := &service{
		items: make(map[string]map[string]*db.Item),
	}

	return s, nil
}

type service struct {
	mu sync.RWMutex

	items map[string]map[string]*db.Item
}

func (s *service) New(i *db.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.items[i.GuildID]
	if !ok {
		s.items[i.GuildID] = make(map[string]*db.Item)
	}

	s.items[i.GuildID][i.Message] = i

	return nil
}

func (s *service) Get(guildID, n string) (*db.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	m, ok := s.items[guildID]
	if !ok {
		return nil, db.ErrItemNotFound
	}

	i, ok := m[n]
	if !ok {
		return nil, db.ErrItemNotFound
	}

	println(guildID, n)

	return i, nil
}

func (s *service) Remove(guildID, n string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.items[guildID]
	if !ok {
		return nil
	}
	delete(s.items[guildID], n)

	return nil
}

func (s *service) All() ([]*db.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var is []*db.Item

	for _, m := range s.items {
		for _, i := range m {
			is = append(is, i)
		}
	}

	return is, nil
}

func (s *service) Ping() error {
	return nil
}
