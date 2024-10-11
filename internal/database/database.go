package database

import (
	"database/sql"
	"os"

	sqliteGo "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	CustomDriverName = "sqlite3_extended"
	File             = "fork_data/forkman.db"
)

func New(log *zerolog.Logger) *gorm.DB {
	sql.Register(CustomDriverName,
		&sqliteGo.SQLiteDriver{
			ConnectHook: func(conn *sqliteGo.SQLiteConn) error {
				err := conn.RegisterFunc(
					"gen_random_uuid",
					func(arguments ...interface{}) (string, error) {
						return uuid.NewV4().String(), nil // Return a string value.
					},
					true,
				)
				return err
			},
		},
	)

	// Make fork_data directory if it doesn't exist
	err := os.MkdirAll("./fork_data", os.ModePerm)
	if err != nil {
		panic(err)
	}

	// Open a new connection to the database
	log.Info().Msg("Opening database connection")
	conn, err := sql.Open(CustomDriverName, File)
	if err != nil {
		panic(err)
	}

	// Register the custom driver
	db, err := gorm.Open(sqlite.Dialector{
		DriverName: CustomDriverName,
		DSN:        File,
		Conn:       conn,
	}, &gorm.Config{
		// Logger:                   logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction:   true,
		DisableNestedTransaction: true,
	})

	// Catalog all models
	models := []interface{}{
		&User{},
		&Module{},
		&Guild{},
		&Email{},
	}

	// Auto migrate the database
	log.Info().Msg("Migrating database models")
	err = db.AutoMigrate(models...)
	if err != nil {
		panic(err)
	}

	return db
}
