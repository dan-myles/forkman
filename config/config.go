package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	DiscordAppID        string `env:"DISCORD_APP_ID required"`
	DiscordClientID     string `env:"DISCORD_CLIENT_ID required"`
	DiscordClientSecret string `env:"DISCORD_CLIENT_SECRET required"`
	DiscordBotToken     string `env:"DISCORD_BOT_TOKEN required"`
	DiscordDevGuildID   string `env:"DISCORD_DEV_GUILD_ID required"`
	DiscordOwnerID      string `env:"DISCORD_OWNER_ID required"`
	LogLevel            string `env:"LOG_LEVEL default=info"`
	GoEnv               string `env:"GO_ENV default=development"`
}

func New() *Config {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		panic(err)
	}

	return &cfg
}
