package postgres

import (
	"database/sql"
	"fmt"

	"github.com/fvdveen/mu2/config"
	"github.com/fvdveen/mu2/db"

	// Register the postgres driver
	_ "github.com/lib/pq"
)

var _ = db.Register("postgres", POSTGRES)

type postgres struct {
	db *sql.DB
}

// POSTGRES creates a new db.Service with a postgres backend
func POSTGRES(conf config.Database) (db.Service, error) {
	db := &postgres{}
	connStr := fmt.Sprintf("user=%s password=%s host=%s sslmode=%s dbname=mu2 port=5432", conf.User, conf.Password, conf.Host, conf.SSL)

	sqlDb, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	_, err = sqlDb.Exec(`CREATE TABLE IF NOT EXISTS commands (
	id SERIAL NOT NULL PRIMARY KEY,
	guild_id TEXT NOT NULL,
	name TEXT NOT NULL,
	response TEXT NOT NULL)`)
	if err != nil {
		return nil, err
	}

	db.db = sqlDb

	return db, nil
}

func (p *postgres) Command(gID, n string) (*db.Command, error) {
	c := &db.Command{}

	stmt, err := p.db.Prepare("SELECT id, guild_id, name, response FROM commands WHERE guild_id=$1 AND name=$2")
	if err != nil {
		return nil, err
	}
	if err := stmt.QueryRow(gID, n).Scan(&c.ID, &c.GID, &c.Name, &c.Response); err != nil {
		if err == sql.ErrNoRows {
			return nil, db.ErrNoCommand
		}
		return nil, err
	}

	return c, nil
}

func (p *postgres) AddCommand(c *db.Command) error {
	stmt, err := p.db.Prepare(`INSERT INTO commands (guild_id, name, response)
	VALUES ($1, $2, $3)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(c.GID, c.Name, c.Response)
	return err
}

func (p *postgres) RemoveCommand(gID, n string) error {
	stmt, err := p.db.Prepare("DELETE FROM commands WHERE guild_id=$1 AND name=$2")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(gID, n)
	return err
}
