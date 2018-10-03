package config

import (
	"io"
)

// Config holds all config values
type Config struct {
	Bot      Bot      `mapstructure:"bot" json:"bot"`
	Log      Log      `mapstructure:"log" json:"log"`
	Database Database `mapstructure:"database" json:"database"`
}

type Bot struct {
	Discord  Discord  `mapstructure:"discord" json:"discord"`
	Prefix   string   `mapstructure:"prefix" json:"prefix"`
	Commands []string `mapstructure:"commands" json:"commands"`
}

// Discord holds all config values for discord
type Discord struct {
	Token string `mapstructure:"token" json:"token"`
}

// Log holds all config values for the logger
type Log struct {
	Discord struct {
		Level   string `mapstructure:"level" json:"level"`
		WebHook string `mapstructure:"webhook" json:"webhook"`
	} `mapstructure:"discord" json:"discord"`

	Level string `mapstructure:"level" json:"level"`
}

// Database holds all config values for the database
type Database struct {
	Host     string `mapstructure:"host" json:"host"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	SSL      string `mapstructure:"ssl" json:"ssl"`
	Type     string `mapstructure:"type" json:"type"`
}

type Watcher interface {
	Watch() <-chan *Config
}

type Provider interface {
	Watcher
	io.Closer
}
