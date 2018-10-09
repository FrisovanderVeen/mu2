package wrapper

import (
	"context"
	"sync"

	"github.com/fvdveen/mu2/services/encode"
)

// Service is a wrapper around the encode service
type Service interface {
	encode.Service
	SetService(s encode.Service) error
}

type service struct {
	s  encode.Service
	mu sync.RWMutex
}

// New creates a new wrapper for the encode service
func New(s encode.Service) Service {
	return &service{
		s: s,
	}
}

func (s *service) Encode(ctx context.Context, url string) (encode.OpusReader, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.s.Encode(ctx, url)
}

func (s *service) SetService(srv encode.Service) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s = srv
	return nil
}
