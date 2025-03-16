package config

import (
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type ForkConfig struct {
	// Server Settings
	DiscordAppID          string        `env:"DISCORD_APP_ID,required,notEmpty"`
	DiscordClientID       string        `env:"DISCORD_CLIENT_ID,required,notEmpty"`
	DiscordClientSecret   string        `env:"DISCORD_CLIENT_SECRET,required,notEmpty"`
	DiscordBotToken       string        `env:"DISCORD_BOT_TOKEN,required,notEmpty"`
	DiscordOwnerID        string        `env:"DISCORD_OWNER_ID,required,notEmpty"`
	ServerPort            int           `env:"SERVER_PORT,required,notEmpty"`
	ServerTimeoutRead     time.Duration `env:"SERVER_TIMEOUT_READ,required,notEmpty"`
	ServerTimeoutWrite    time.Duration `env:"SERVER_TIMEOUT_WRITE,required,notEmpty"`
	ServerTimeoutIdle     time.Duration `env:"SERVER_TIMEOUT_IDLE,required,notEmpty"`
	ServerAuthSecret      string        `env:"SERVER_AUTH_SECRET,required,notEmpty"`
	ServerAuthExpiry      time.Duration `env:"SERVER_AUTH_EXPIRY,required,notEmpty"`
	ServerAuthCallbackURI string        `env:"SERVER_AUTH_CALLBACK_URI,required,notEmpty"`
	LogLevel              string        `env:"LOG_LEVEL,required,notEmpty"`
	GoEnv                 string        `env:"GO_ENV,required,notEmpty"`

	// AWS
	AWS_ACCESS_KEY_ID     string `env:"AWS_ACCESS_KEY_ID,required,notEmpty"`
	AWS_SECRET_ACCESS_KEY string `env:"AWS_SECRET_ACCESS_KEY,required,notEmpty"`
	AWS_REGION            string `env:"AWS_REGION,required,notEmpty"`
	AWS_BEDROCK_KBI       string `env:"AWS_BEDROCK_KBI,required,notEmpty"` // Knowledge Base ID

	// Verification Settings
	LogChannelID string `env:"LOG_CHANNEL_ID,required,notEmpty"`
	RoleToRemove string `env:"ROLE_TO_REMOVE,required,notEmpty"`
	RoleToAdd    string `env:"ROLE_TO_ADD,required,notEmpty"`

	// QNA Settings
	FORUM_CHANNEL_ID string `env:"FORUM_CHANNEL_ID,required,notEmpty"`
}

func New() *ForkConfig {
	if os.Getenv("GO_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, relying on external env vars...")
		}
	}

	cfg, err := env.ParseAs[ForkConfig]()
	if err != nil {
		panic(err)
	}

	return &cfg
}
