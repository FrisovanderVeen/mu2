package config

import (
	"os"
)

// Config holds all config values
type Config struct {
	Discord  Discord
	Log      Log
	Database Database
}

// Discord holds all config values for discord
type Discord struct {
	Token  string
	Prefix string
}

// Log holds all config values for the logger
type Log struct {
	DiscordWebHook string
	DiscordLevel   string

	Level string
}

// Database holds all config values for the database
type Database struct {
	Host     string
	User     string
	Password string
	SSL      string
	Type     string
}

// Load loads the config values from environment variables
func Load() *Config {
	c := &Config{
		Discord: Discord{
			Token:  os.Getenv("DISCORD_TOKEN"),
			Prefix: os.Getenv("DISCORD_PREFIX"),
		},
		Log: Log{
			DiscordWebHook: os.Getenv("LOG_WEBHOOK_DISCORD"),
			DiscordLevel:   os.Getenv("LOG_LEVEL_DISCORD"),
			Level:          os.Getenv("LOG_LEVEL"),
		},
		Database: Database{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
			SSL:      os.Getenv("DB_SSL"),
			Type:     os.Getenv("DB_TYPE"),
		},
	}

	return c
}

// Defaults sets the unset values to their defaults
func (c *Config) Defaults() {
	if c.Discord.Prefix == "" {
		c.Discord.Prefix = "$"
	}
	if c.Database.Type == "" {
		c.Database.Type = "postgres"
	}
}
