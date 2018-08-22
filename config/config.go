package config

import (
	"gopkg.in/urfave/cli.v2"
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

// Load loads the values of hte context into a Config
func Load(c *cli.Context) *Config {
	conf := &Config{
		Discord{},
		Log{},
		Database{},
	}

	if token := c.String("token"); token != "" {
		conf.Discord.Token = token
	}
	if prefix := c.String("prefix"); prefix != "" {
		conf.Discord.Prefix = prefix
	}
	if lvl := c.String("log-level"); lvl != "" {
		conf.Log.Level = lvl
	}
	if hook := c.String("discord-webhook"); hook != "" {
		conf.Log.DiscordWebHook = hook
	}
	if dlvl := c.String("discord-log-level"); dlvl != "" {
		conf.Log.DiscordLevel = dlvl
	}
	if dbHost := c.String("db-host"); dbHost != "" {
		conf.Database.Host = dbHost
	}
	if dbUser := c.String("db-user"); dbUser != "" {
		conf.Database.User = dbUser
	}
	if dbPass := c.String("db-password"); dbPass != "" {
		conf.Database.Password = dbPass
	}
	if dbSSL := c.String("db-ssl"); dbSSL != "" {
		conf.Database.SSL = dbSSL
	}
	if dbType := c.String("db-type"); dbType != "" {
		conf.Database.Type = dbType
	}

	return conf
}
