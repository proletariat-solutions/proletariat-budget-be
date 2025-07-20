package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

// MigrationConfig holds configuration for database migrations
type MigrationConfig struct {
	MigrationsPath string
	DBName         string
	DBUser         string
	DBPassword     string
	DBHost         string
	DBPort         string
}

// RunMigrations runs all pending migrations
func RunMigrations(config MigrationConfig) error {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?multiStatements=true",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	db, err := sql.Open(
		"mysql",
		dsn,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to open database connection: %w",
			err,
		)
	}
	defer db.Close()

	maxOpenConnections := 10
	maxIdleConnections := 5
	fiveMinutes := time.Minute * 5

	// Set connection pool settings
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetConnMaxLifetime(fiveMinutes)

	// Ping database to verify connection
	if errPing := db.PingContext(context.Background()); errPing != nil {
		return fmt.Errorf(
			"failed to ping database: %w",
			errPing,
		)
	}

	// Create migration driver
	driver, err := mysql.WithInstance(
		db,
		&mysql.Config{},
	)
	if err != nil {
		return fmt.Errorf(
			"failed to create migration driver: %w",
			err,
		)
	}

	// Create migration instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+config.MigrationsPath,
		"mysql",
		driver,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to create migration instance: %w",
			err,
		)
	}

	// Run migrations
	log.Info().Msg("Running database migrations...")
	if errUp := m.Up(); errUp != nil && !errors.Is(
		errUp,
		migrate.ErrNoChange,
	) {
		return fmt.Errorf(
			"failed to run migrations: %w",
			errUp,
		)
	}

	// Get migration version
	version, dirty, err := m.Version()
	if err != nil && !errors.Is(
		err,
		migrate.ErrNilVersion,
	) {
		return fmt.Errorf(
			"failed to get migration version: %w",
			err,
		)
	}

	log.Info().
		Uint(
			"version",
			version,
		).
		Bool(
			"dirty",
			dirty,
		).
		Msg("Database migration completed successfully")

	return nil
}
