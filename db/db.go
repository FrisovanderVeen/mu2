package db

import (
	"errors"

	"github.com/fvdveen/mu2/config"
)

var dbs = map[string]DBFunc{}

var (
	// ErrDBNotFound is used when the given database type is not found
	ErrDBNotFound = errors.New("db not found")
	// ErrNoCommand is used when the given command is not found
	ErrNoCommand = errors.New("command doesn't exist")
)

// DBFunc creates a new database
type DBFunc func(config.Database) (Service, error)

// Service is a service that holds learned commands
type Service interface {
	Command(gID, n string) (*Command, error)
	AddCommand(*Command) error
	RemoveCommand(gID, n string) error
}

// Command represents a command in the database
type Command struct {
	ID       int
	GID      string
	Name     string
	Response string
}

// Register adds a new DBFunc
func Register(n string, db DBFunc) DBFunc {
	dbs[n] = db

	return db
}

// Get creates a new database
func Get(c config.Database) (Service, error) {
	dbFunc, ok := dbs[c.Type]
	if !ok {
		return nil, ErrDBNotFound
	}
	return dbFunc(c)
}
