package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type DiscordConfig struct {
	AppID          string `env:"DISCORD_APP_ID,required,notEmpty"`
	ClientID       string `env:"DISCORD_CLIENT_ID,required,notEmpty"`
	ClientSecret   string `env:"DISCORD_CLIENT_SECRET,required,notEmpty"`
	BotToken       string `env:"DISCORD_BOT_TOKEN,required,notEmpty"`
	OwnerID        string `env:"DISCORD_OWNER_ID,required,notEmpty"`
	LogChannelID   string `env:"LOG_CHANNEL_ID,required,notEmpty"`
	RoleToRemove   string `env:"ROLE_TO_REMOVE,required,notEmpty"`
	RoleToAdd      string `env:"ROLE_TO_ADD,required,notEmpty"`
	FORUMChannelID string `env:"FORUM_CHANNEL_ID,required,notEmpty"`
}

type AWSConfig struct {
	AWSAccessKeyID     string `env:"AWS_ACCESS_KEY_ID,required,notEmpty"`
	AWSSecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY,required,notEmpty"`
	AWSRegion          string `env:"AWS_REGION,required,notEmpty"`
	AWSBedrockKBI      string `env:"AWS_BEDROCK_KBI,required,notEmpty"` // Knowledge Base ID
}

type ServerConfig struct {
	Port         int           `env:"SERVER_PORT,required,notEmpty"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required,notEmpty"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required,notEmpty"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required,notEmpty"`
	AuthSecret   string        `env:"SERVER_AUTH_SECRET,required,notEmpty"`
	AuthExpiry   time.Duration `env:"SERVER_AUTH_EXPIRY,required,notEmpty"`
}

type SentinelConfig struct {
	// Server Config
	ServerConfig *ServerConfig `env:"-"`

	// AWS Config (optional)
	AWSEnabled bool       `env:"AWS_ENABLED`
	AWSConfig  *AWSConfig `env:"-"`

	// Discord Config
	DiscordConfig *DiscordConfig `env:"-"`

	// General Config
	LogLevel string `env:"LOG_LEVEL,required,notEmpty"`
	GoEnv    string `env:"GO_ENV,required,notEmpty"`
}

func New() *SentinelConfig {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	cfg := SentinelConfig{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	// Load Server Config
	serverCfg := ServerConfig{}
	if err := env.Parse(&serverCfg); err != nil {
		panic(err)
	}
	cfg.ServerConfig = &serverCfg

	// Load AWS Config if AWS is enabled
	if cfg.AWSEnabled {
		awsCfg := AWSConfig{}
		if err := env.Parse(&awsCfg); err != nil {
			panic(err)
		}
		cfg.AWSConfig = &awsCfg
	}

	// Load Discord Config
	discordCfg := DiscordConfig{}
	if err := env.Parse(&discordCfg); err != nil {
		panic(err)
	}
	cfg.DiscordConfig = &discordCfg

	return &cfg
}
