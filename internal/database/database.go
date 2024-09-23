package database

import (
	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New(log *zerolog.Logger) *gorm.DB {
	// Open a new connection to the database
	log.Info().Msg("Opening database connection")
	config := &gorm.Config{}
	db, err := gorm.Open(sqlite.Open("fork_data/forkman.db?_foreign_keys=on"), config)
	if err != nil {
		panic(err)
	}

	// Catalog all models
	models := []interface{}{
		&User{},
		&Module{},
		&Guild{},
	}

	// Auto migrate the database
	log.Info().Msg("Migrating database models")
	err = db.AutoMigrate(models...)
	if err != nil {
		panic(err)
	}

	return db
}
