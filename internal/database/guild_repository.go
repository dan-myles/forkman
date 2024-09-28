package database

import (
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

const iconSize = "128"

type GuildRepository struct {
	db *gorm.DB
}

func NewGuildRepository(db *gorm.DB) *GuildRepository {
	return &GuildRepository{
		db: db,
	}
}

func (r *GuildRepository) UpdateGuild(guild *discordgo.Guild) (*Guild, error) {
	g := &Guild{}
	result := r.db.First(g, "snowflake = ?", guild.ID)
	if result.Error != nil {
		return nil, result.Error
	}

	g.Name = guild.Name
	g.IconUrl = guild.IconURL(iconSize)
	g.OwnerID = guild.OwnerID

	err := r.db.Save(g).Error
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (r *GuildRepository) CreateGuild(guild *discordgo.Guild) (*Guild, error) {
	g := &Guild{
		Snowflake: guild.ID,
		Name:      guild.Name,
		IconUrl:   guild.IconURL(iconSize),
		OwnerID:   guild.OwnerID,
	}

	if err := r.db.Create(g).Error; err != nil {
		return nil, err
	}

	return g, nil
}

func (r *GuildRepository) ReadGuild(guildSnowflake string) (*Guild, error) {
	g := &Guild{}
	err := r.db.First(g, "snowflake = ?", guildSnowflake).Error
	if err != nil {
		return nil, err
	}

	return g, nil
}
