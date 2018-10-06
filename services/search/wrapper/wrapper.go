package wrapper

import (
	"context"
	"sync"

	"github.com/fvdveen/mu2/services/search"
)

// Service is a wrapper around the search service
type Service interface {
	search.Service
	SetService(s search.Service) error
}

type service struct {
	s  search.Service
	mu sync.RWMutex
}

// New creates a new wrapper for the search service
func New(s search.Service) Service {
	return &service{
		s: s,
	}
}

func (s *service) Search(ctx context.Context, n string) (*search.Video, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.s.Search(ctx, n)
}

func (s *service) SetService(srv search.Service) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s = srv
	return nil
}
