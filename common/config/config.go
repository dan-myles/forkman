package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type SentinelConfig struct {
	DiscordAppID        string        `env:"DISCORD_APP_ID,required,notEmpty"`
	DiscordClientID     string        `env:"DISCORD_CLIENT_ID,required,notEmpty"`
	DiscordClientSecret string        `env:"DISCORD_CLIENT_SECRET,required,notEmpty"`
	DiscordBotToken     string        `env:"DISCORD_BOT_TOKEN,required,notEmpty"`
	DiscordDevGuildID   string        `env:"DISCORD_DEV_GUILD_ID,required,notEmpty"`
	DiscordOwnerID      string        `env:"DISCORD_OWNER_ID,required,notEmpty"`
	ServerPort          int           `env:"SERVER_PORT,required,notEmpty"`
	ServerTimeoutRead   time.Duration `env:"SERVER_TIMEOUT_READ,required,notEmpty"`
	ServerTimeoutWrite  time.Duration `env:"SERVER_TIMEOUT_WRITE,required,notEmpty"`
	ServerTimeoutIdle   time.Duration `env:"SERVER_TIMEOUT_IDLE,required,notEmpty"`
	ServerAuthSecret    string        `env:"SERVER_AUTH_SECRET,required,notEmpty"`
	ServerAuthExpiry    time.Duration `env:"SERVER_AUTH_EXPIRY,required,notEmpty"`
	LogLevel            string        `env:"LOG_LEVEL,required,notEmpty"`
	GoEnv               string        `env:"GO_ENV,required,notEmpty"`
	ResendAPIKey        string        `env:"RESEND_API_KEY,required,notEmpty"`
	LogChannelID        string        `env:"LOG_CHANNEL_ID,required,notEmpty"`
	RoleToRemove        string        `env:"ROLE_TO_REMOVE,required,notEmpty"`
	RoleToAdd           string        `env:"ROLE_TO_ADD,required,notEmpty"`
}

func New() *SentinelConfig {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	cfg, err := env.ParseAs[SentinelConfig]()
	if err != nil {
		panic(err)
	}

	return &cfg
}
