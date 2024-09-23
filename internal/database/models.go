package database

import (
	"time"
)

type User struct {
	ID               string `gorm:"primarykey"`
	DiscordID        string
	DiscordUsername  string
	DiscordAvatarURL string
	DiscordEmail     string
	LastLogin        time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Guild struct {
	ID      string   `gorm:"primarykey"`
	Modules []Module `gorm:"foreignKey:GuildID;references:ID"`
}

type Module struct {
	ID          string `gorm:"primarykey"`
	GuildID     uint
	Name        string
	Description string
	Enabled     bool `gorm:"default:false"`
}
