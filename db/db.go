package db

import (
	"errors"

	"github.com/fvdveen/mu2-config"
)

var fs = make(map[string]Factory)

var (
	// ErrDBNotFound is used when a db is not found
	ErrDBNotFound = errors.New("db type not found")

	// ErrItemNotFound is used when an item is not found in the db
	ErrItemNotFound = errors.New("item not found in db")
)

// Factory creates a service
type Factory func(config.Database) (Service, error)

// Service holds Items
type Service interface {
	New(*Item) error
	// Get takes in the guildID and the name of the command and returns the stored item
	Get(string, string) (*Item, error)
	// Remove takes in the guildID and the name of the command and removes it
	Remove(string, string) error
	All() ([]*Item, error)
	// Ping is a health check for the db
	// if it fails it returns a non-nil error
	Ping() error
}

// Item is a key-value mapped item
type Item struct {
	GuildID, Message, Response string
}

// Register registers a factory
func Register(name string, f Factory) Factory {
	fs[name] = f
	return f
}

// Get creates a service based on the given configuration
func Get(conf config.Database) (Service, error) {
	f, ok := fs[conf.Type]
	if !ok {
		return nil, ErrDBNotFound
	}

	return f(conf)
}
