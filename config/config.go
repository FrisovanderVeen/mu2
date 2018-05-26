package config

// Config holds all configuration values
type Config struct {
	Bot Bot `mapstructure:"bot"`
}

// Bot holds all configuration values for the discord bot
type Bot struct {
	Token      string `mapstructure:"token"`
	Prefix     string `mapstructure:"prefix"`
	InviteLink string `mapstructure:"invite-link"`
}
