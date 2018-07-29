package bot

import (
	"github.com/fvdveen/mu2/config"
	"github.com/fvdveen/mu2/db"
)

// OptionFunc is a option for New
type OptionFunc func(*Bot)

// WithConfig sets the bot's config
func WithConfig(c config.Config) OptionFunc {
	return func(b *Bot) {
		b.conf = c
	}
}

// WithDB sets the bot's database
func WithDB(db db.Service) OptionFunc {
	return func(b *Bot) {
		b.db = db
	}
}
