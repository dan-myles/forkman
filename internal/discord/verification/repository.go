package verification

import (
	"github.com/avvo-na/forkman/internal/database"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateModule(mod *database.Module) (*database.Module, error) {
	result := r.db.Create(mod)
	if result.Error != nil {
		return nil, result.Error
	}

	return mod, nil
}

func (r *Repository) ReadModule(guildSnowflake string) (*database.Module, error) {
	mod := &database.Module{}
	result := r.db.First(mod, "name = ? AND guild_snowflake = ?", name, guildSnowflake)
	if result.Error != nil {
		return nil, result.Error
	}

	return mod, nil
}

func (r *Repository) UpdateModule(mod *database.Module) (*database.Module, error) {
	m := &database.Module{}
	result := r.db.First(m, "name = ? AND guild_snowflake = ?", name, mod.GuildSnowflake)
	if result.Error != nil {
		return nil, result.Error
	}

	m.Enabled = mod.Enabled
	m.Config = mod.Config
	m.Commands = mod.Commands

	err := r.db.Save(m).Error
	if err != nil {
		return nil, err
	}

	return nil, nil
}
