package database

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID             string
	DiscordID        string
	DiscordUsername  string
	DiscordAvatarURL string
	DiscordEmail     string
	LastLogin        time.Time
}
