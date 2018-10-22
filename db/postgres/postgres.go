package postgres

import (
	"database/sql"
	"fmt"

	"github.com/fvdveen/mu2-config"
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
	connStr := fmt.Sprintf("user=%s password=%s host=%s sslmode=%s dbname=mu2 port=%s", conf.User, conf.Password, conf.Host, conf.SSL, conf.Port)

	sqlDb, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	_, err = sqlDb.Exec(`CREATE TABLE IF NOT EXISTS commands (
	id SERIAL NOT NULL PRIMARY KEY,
	guild_id TEXT NOT NULL,
	message TEXT NOT NULL,
	response TEXT NOT NULL)`)
	if err != nil {
		return nil, err
	}

	db.db = sqlDb

	return db, nil
}

func (p *postgres) Get(gID, n string) (*db.Item, error) {
	i := &db.Item{}

	stmt, err := p.db.Prepare("SELECT guild_id, message, response FROM commands WHERE guild_id=$1 AND message=$2")
	if err != nil {
		return nil, err
	}
	if err := stmt.QueryRow(gID, n).Scan(&i.GuildID, &i.Message, &i.Response); err != nil {
		if err == sql.ErrNoRows {
			return nil, db.ErrItemNotFound
		}
		return nil, err
	}

	return i, nil
}

func (p *postgres) New(c *db.Item) error {
	stmt, err := p.db.Prepare(`INSERT INTO commands (guild_id, message, response)
	VALUES ($1, $2, $3)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(c.GuildID, c.Message, c.Response)
	return err
}

func (p *postgres) Remove(gID, n string) error {
	stmt, err := p.db.Prepare("DELETE FROM commands WHERE guild_id=$1 AND message=$2")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(gID, n)
	return err
}

func (p *postgres) All() ([]*db.Item, error) {
	stmt, err := p.db.Prepare("SELECT guild_id, message, response FROM commands")
	if err != nil {
		return nil, fmt.Errorf("prepare statement: %v", err)
	}

	rs, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("exec statement: %v", err)
	}

	var is []*db.Item
	for rs.Next() {
		i := &db.Item{}
		if err := rs.Scan(&i.GuildID, &i.Message, &i.Response); err != nil {
			return nil, fmt.Errorf("scan row: %v", err)
		}
		is = append(is, i)
	}
	if err := rs.Err(); err != nil {
		return nil, fmt.Errorf("reading rows: %v", err)
	}

	return is, nil
}

func (p *postgres) Ping() error {
	return p.db.Ping()
}
