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
	m := &database.Module{}
	result := r.db.First(m, "name = ? AND guild_snowflake = ?", name, guildSnowflake)
	if result.Error != nil {
		return nil, result.Error
	}

	return m, nil
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

	return m, nil
}

func (r *Repository) ReadEmail(guildSnowflake string, userSnowflake string) (*database.Email, error) {
	e := &database.Email{}
	result := r.db.First(e, "guild_snowflake = ? AND user_snowflake = ?", guildSnowflake, userSnowflake)
	if result.Error != nil {
		return nil, result.Error
	}

	return e, nil
}

func (r *Repository) UpdateEmail(email *database.Email) (*database.Email, error) {
	e := &database.Email{}
	result := r.db.First(e, "guild_snowflake = ? AND user_snowflake = ?", email.GuildSnowflake, email.UserSnowflake)
	if result.Error != nil {
		return nil, result.Error
	}

	e.Address = email.Address
	e.Code = email.Code
	e.IsVerified = email.IsVerified

	err := r.db.Save(e).Error
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (r *Repository) UpsertEmail(email *database.Email) (*database.Email, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		existingEmail := &database.Email{}
		result := tx.Where("user_snowflake = ? AND guild_snowflake = ?", email.UserSnowflake, email.GuildSnowflake).First(existingEmail)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				if err := tx.Create(email).Error; err != nil {
					return err
				}
				return nil
			}

			return result.Error
		}

		existingEmail.Address = email.Address
		existingEmail.Code = email.Code
		existingEmail.IsVerified = email.IsVerified
		if err := tx.Save(existingEmail).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return email, nil
}
