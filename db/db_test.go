package db_test

import (
	"testing"

	"github.com/fvdveen/mu2/config"
	"github.com/fvdveen/mu2/db"
)

func TestDB(t *testing.T) {
	_ = db.Register("test", func(config.Database) (db.Service, error) {
		return nil, nil
	})

	s, err := db.Get(config.Database{
		Type: "test",
	})
	if err != nil {
		t.Errorf("Expcted: %v, got: %v", nil, err)
	}
	if s != nil {
		t.Errorf("Expcted: %v, got: %v", nil, nil)
	}

	s, err = db.Get(config.Database{
		Type: "test1",
	})
	if err != db.ErrDBNotFound {
		t.Errorf("Expcted: %v, got: %v", nil, err)
	}
	if s != nil {
		t.Errorf("Expcted: %v, got: %v", nil, nil)
	}
}
