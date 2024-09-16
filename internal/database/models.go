package database

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID               string
	DiscordID            string
	DiscordUsername      string
	DiscordDiscriminator string
	DiscordAvatarURL     string
	DiscordEmail         string
	LastLogin            time.Time
}
