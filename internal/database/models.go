package database

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/datatypes"
)

type User struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:(gen_random_uuid())"`
	DiscordSnowflake string    `gorm:"uniqueIndex"`
	DiscordUsername  string
	DiscordAvatarURL string
	DiscordEmail     string `gorm:"uniqueIndex"`
	LastLogin        time.Time
	CreatedAt        time.Time // Managed by GORM
	UpdatedAt        time.Time // Managed by GORM
}

type Guild struct {
	Snowflake  string    `gorm:"primaryKey;unique"`
	Modules    []Module  `gorm:"foreignKey:GuildSnowflake;references:Snowflake;constraint:OnDelete:CASCADE"`
	Admins     []User    `gorm:"many2many:guild_admins;"`
	AdminRoles []string  `gorm:"type:text[]"`
	CreatedAt  time.Time // Managed by GORM
	UpdatedAt  time.Time // Managed by GORM
}

type Module struct {
	ID             uint   `gorm:"primarykey;autoIncrement"`
	GuildSnowflake string `gorm:"index"`
	Name           string
	Description    string
	Enabled        bool `gorm:"default:false"`
	State          datatypes.JSON
	CreatedAt      time.Time // Managed by GORM
	UpdatedAt      time.Time // Managed by GORM
}
