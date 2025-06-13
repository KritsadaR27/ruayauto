package database

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/stdlib"
)

// MigrationConfig holds migration configuration
type MigrationConfig struct {
	MigrationsPath string
	DatabaseName   string
}

// RunMigrations runs database migrations from the specified directory
func (db *DB) RunMigrations(ctx context.Context, config *MigrationConfig) error {
	log.Printf("Starting database migrations...")

	// Get a standard database connection for migrations
	stdDB := stdlib.OpenDBFromPool(db.Pool)
	defer stdDB.Close()

	// Create postgres driver instance
	driver, err := postgres.WithInstance(stdDB, &postgres.Config{
		DatabaseName: config.DatabaseName,
	})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Get absolute path to migrations
	migrationsPath, err := filepath.Abs(config.MigrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute migrations path: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		config.DatabaseName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Run migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Printf("No new migrations to apply")
	} else {
		log.Printf("Migrations applied successfully")
	}

	// Get current version
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if err == migrate.ErrNilVersion {
		log.Printf("Database schema is at version: nil (no migrations applied)")
	} else {
		log.Printf("Database schema is at version: %d (dirty: %t)", version, dirty)
	}

	return nil
}

// CreateMigration creates a new migration file pair (up/down) in the migrations directory
func CreateMigration(migrationsPath, name string) error {
	// This function would create migration files
	// For now, we'll just return instructions
	log.Printf("To create a new migration, run:")
	log.Printf("migrate create -ext sql -dir %s -seq %s", migrationsPath, name)
	return nil
}

// RollbackMigration rolls back the last migration
func (db *DB) RollbackMigration(ctx context.Context, config *MigrationConfig) error {
	log.Printf("Rolling back last migration...")

	// Get a standard database connection for migrations
	stdDB := stdlib.OpenDBFromPool(db.Pool)
	defer stdDB.Close()

	// Create postgres driver instance
	driver, err := postgres.WithInstance(stdDB, &postgres.Config{
		DatabaseName: config.DatabaseName,
	})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Get absolute path to migrations
	migrationsPath, err := filepath.Abs(config.MigrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute migrations path: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		config.DatabaseName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Roll back one step
	err = m.Steps(-1)
	if err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Printf("Migration rolled back successfully")
	return nil
}

// MigrationStatus returns the current migration status
func (db *DB) MigrationStatus(ctx context.Context, config *MigrationConfig) (uint, bool, error) {
	// Get a standard database connection for migrations
	stdDB := stdlib.OpenDBFromPool(db.Pool)
	defer stdDB.Close()

	// Create postgres driver instance
	driver, err := postgres.WithInstance(stdDB, &postgres.Config{
		DatabaseName: config.DatabaseName,
	})
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Get absolute path to migrations
	migrationsPath, err := filepath.Abs(config.MigrationsPath)
	if err != nil {
		return 0, false, fmt.Errorf("failed to get absolute migrations path: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		config.DatabaseName,
		driver,
	)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	// Get current version
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	if err == migrate.ErrNilVersion {
		return 0, false, nil
	}

	return version, dirty, nil
}
