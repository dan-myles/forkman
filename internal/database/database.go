package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New() *gorm.DB {
	// Open a new connection to the database
	config := &gorm.Config{}
	db, err := gorm.Open(sqlite.Open("fork_data/forkman.db"), config)
	if err != nil {
		panic(err)
	}

	// Catalog all models
	models := []interface{}{
		&User{},
	}

	// Auto migrate the database
	err = db.AutoMigrate(models...)
	if err != nil {
		panic(err)
	}

	return db
}
