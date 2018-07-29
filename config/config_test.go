package config_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/fvdveen/mu2/config"
)

func TestLoad(t *testing.T) {
	envs := map[string]string{
		"DISCORD_TOKEN":       "haha nice try",
		"DISCORD_PREFIX":      "$",
		"LOG_WEBHOOK_DISCORD": "https://somewebhook.com",
		"LOG_LEVEL_DISCORD":   "none",
		"LOG_LEVEL":           "",
		"DB_HOST":             "DB_HOST",
		"DB_USER":             "DB_USER",
		"DB_PASS":             "DB_PASS",
	}

	for k, v := range envs {
		if err := os.Setenv(k, v); err != nil {
			t.Error(err)
		}
	}

	expected := config.Config{
		Discord: config.Discord{
			Token:  "haha nice try",
			Prefix: "$",
		},
		Log: config.Log{
			DiscordWebHook: "https://somewebhook.com",
			DiscordLevel:   "none",
			Level:          "",
		},
		Database: config.Database{
			Host:     "DB_HOST",
			User:     "DB_USER",
			Password: "DB_PASS",
		},
	}

	conf := config.Load()
	if !reflect.DeepEqual(*conf, expected) {
		t.Errorf("Got: %+v, expected: %+v", conf, expected)
	}
}

func TestDefaults(t *testing.T) {
	expected := &config.Config{
		Discord: config.Discord{
			Prefix: "$",
		},
		Database: config.Database{
			Type: "postgres",
		},
	}

	c := &config.Config{}
	c.Defaults()

	if !reflect.DeepEqual(expected, c) {
		t.Errorf("Expected: %+v, got: %+v", expected, c)
	}
}
