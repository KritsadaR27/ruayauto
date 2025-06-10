package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"ruayautomsg/internal/config"
	"ruayautomsg/internal/database"
)

func main() {
	var (
		testConnection = flag.Bool("test", false, "Test database connection")
		setup         = flag.Bool("setup", false, "Setup database with migrations")
		verbose       = flag.Bool("verbose", false, "Verbose output")
	)
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if *verbose {
		fmt.Printf("Database URL: %s\n", cfg.GetDatabaseURL())
		fmt.Printf("Migrations Path: ./internal/migrations\n")
	}

	switch {
	case *testConnection:
		if err := testDatabaseConnection(cfg, *verbose); err != nil {
			log.Fatalf("Database connection test failed: %v", err)
		}
		fmt.Println("✅ Database connection test passed!")

	case *setup:
		if err := setupDatabase(cfg, *verbose); err != nil {
			log.Fatalf("Database setup failed: %v", err)
		}
		fmt.Println("✅ Database setup completed successfully!")

	default:
		fmt.Println("Database Test Tool")
		fmt.Println("Usage:")
		fmt.Println("  -test      Test database connection")
		fmt.Println("  -setup     Setup database with migrations")
		fmt.Println("  -verbose   Enable verbose output")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  go run ./cmd/dbtest -test")
		fmt.Println("  go run ./cmd/dbtest -setup")
		fmt.Println("  go run ./cmd/dbtest -test -verbose")
	}
}

func testDatabaseConnection(cfg *config.Config, verbose bool) error {
	if verbose {
		log.Printf("Testing database connection...")
	}
	
	// Create a simple database manager without auto-migrations
	mgr, err := database.NewManagerFromURL(cfg.GetDatabaseURL(), "./internal/migrations", false)
	if err != nil {
		return fmt.Errorf("failed to create database manager: %w", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		mgr.Shutdown(ctx)
	}()

	// Test basic connectivity
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := mgr.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Test health check
	if err := mgr.WaitForHealth(time.Second * 5); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	// Get status
	status := mgr.GetHealthStatus()
	if verbose {
		fmt.Printf("Health status: %t\n", status.IsHealthy)
		fmt.Printf("Pool stats - Total: %d, Idle: %d, Acquired: %d\n", 
			status.PoolStats.TotalConns, 
			status.PoolStats.IdleConns, 
			status.PoolStats.AcquiredConns)
	}

	// Test a simple query
	var result int
	err = mgr.GetDB().QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("test query failed: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("unexpected query result: %d", result)
	}

	if verbose {
		fmt.Printf("Test query successful: %d\n", result)
	}

	return nil
}

func setupDatabase(cfg *config.Config, verbose bool) error {
	if verbose {
		log.Printf("Setting up database...")
	}
	
	// Create database manager with auto-migrations
	mgr, err := database.NewManagerFromURL(cfg.GetDatabaseURL(), "./internal/migrations", true)
	if err != nil {
		return fmt.Errorf("failed to create database manager: %w", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		mgr.Shutdown(ctx)
	}()

	// Wait for health
	if err := mgr.WaitForHealth(time.Second * 30); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	// Check migration status
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	version, dirty, err := mgr.GetDB().MigrationStatus(ctx, &database.MigrationConfig{
		MigrationsPath: "./internal/migrations",
		DatabaseName:   cfg.DBName,
	})
	if err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	if verbose {
		if version == 0 && !dirty {
			fmt.Println("Migration status: No migrations applied yet")
		} else {
			fmt.Printf("Migration status: Version %d (dirty: %t)\n", version, dirty)
		}
	}

	// Test that we can access created tables
	var count int
	err = mgr.GetDB().QueryRow(ctx, `
		SELECT count(*) 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name IN ('keywords', 'message_analytics', 'webhook_logs')
	`).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check tables: %w", err)
	}

	if verbose {
		fmt.Printf("Found %d expected tables in database\n", count)
	}

	if count < 3 {
		return fmt.Errorf("expected at least 3 tables, found %d", count)
	}

	return nil
}
