package wrapper

import (
	"errors"
	"sync"

	"github.com/fvdveen/mu2/db"
)

// Service represents a wrapper around a db.Service
type Service interface {
	db.Service
	// SetService sets the underlying service it will return a non-nil error if the given Service is nil
	SetService(db.Service) error
}

// New returns a new wrapper
func New(s db.Service) Service {
	return &service{
		s: s,
	}
}

type service struct {
	s  db.Service
	mu sync.RWMutex
}

func (s *service) SetService(srv db.Service) error {
	if srv == nil {
		return errors.New("given service is nil")
	}

	s.mu.Lock()
	s.s = srv
	s.mu.Unlock()

	return nil
}

func (s *service) New(i *db.Item) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.s.New(i)
}

func (s *service) Get(gID string, n string) (*db.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.s.Get(gID, n)
}

func (s *service) Remove(gID string, n string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.s.Remove(gID, n)
}

func (s *service) All() ([]*db.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.s.All()
}

func (s *service) Ping() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.s.Ping()
}
