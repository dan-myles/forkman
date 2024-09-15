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
	Session              []Session `gorm:"foreignKey:UserID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Session struct {
	gorm.Model
	SessionID string
	UserID    string
	Token     string
	ExpiresAt time.Time
}
