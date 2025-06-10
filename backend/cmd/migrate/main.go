package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"ruayautomsg/internal/config"
	"ruayautomsg/internal/database"
)

func main() {
	var (
		up        = flag.Bool("up", false, "Run migrations up")
		down      = flag.Bool("down", false, "Run migrations down (rollback)")
		status    = flag.Bool("status", false, "Show migration status")
		create    = flag.String("create", "", "Create new migration with given name")
	)
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create database manager without auto-migrate
	mgr, err := database.NewManagerFromURL(cfg.GetDatabaseURL(), "./internal/migrations", false)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		mgr.Shutdown(ctx)
	}()

	switch {
	case *up:
		log.Println("Running migrations up...")
		if err := mgr.RunMigrations(); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("Migrations completed successfully")

	case *down:
		log.Println("Rolling back migration...")
		if err := mgr.RollbackMigration(); err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
		log.Println("Rollback completed successfully")

	case *status:
		log.Println("Checking migration status...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		
		version, dirty, err := mgr.GetDB().MigrationStatus(ctx, &database.MigrationConfig{
			MigrationsPath: "./internal/migrations",
			DatabaseName:   cfg.DBName,
		})
		if err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
		
		if version == 0 && !dirty {
			fmt.Println("No migrations have been applied yet")
		} else {
			fmt.Printf("Current migration version: %d\n", version)
			fmt.Printf("Dirty state: %t\n", dirty)
		}

	case *create != "":
		log.Printf("Creating migration: %s", *create)
		if err := database.CreateMigration("./internal/migrations", *create); err != nil {
			log.Fatalf("Failed to create migration: %v", err)
		}

	default:
		fmt.Println("Database Migration Tool")
		fmt.Println("Usage:")
		fmt.Println("  -up           Run migrations up")
		fmt.Println("  -down         Run migrations down (rollback)")
		fmt.Println("  -status       Show migration status")
		fmt.Println("  -create NAME  Create new migration")
		os.Exit(1)
	}
}