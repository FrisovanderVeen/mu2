package config

// Config holds all config values
type Config struct {
	Discord  Discord  `mapstructure:"discord"`
	Log      Log      `mapstructure:"log"`
	Database Database `mapstructure:"database"`
}

// Discord holds all config values for discord
type Discord struct {
	Token  string `mapstructure:"token"`
	Prefix string `mapstructure:"prefix"`
}

// Log holds all config values for the logger
type Log struct {
	Discord struct {
		Level   string `mapstructure:"level"`
		WebHook string `mapstructure:"webhook"`
	} `mapstructure:"discord"`

	Level string `mapstructure:"level"`
}

// Database holds all config values for the database
type Database struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSL      string `mapstructure:"ssl"`
	Type     string `mapstructure:"type"`
}
