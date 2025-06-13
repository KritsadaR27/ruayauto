package database

import (
	"context"
	"fmt"
	"log"
	"time"
)

// TestConnection tests the database connection and setup
func TestConnection(databaseURL string) error {
	log.Printf("Testing database connection...")
	
	// Create a simple database manager
	mgr, err := NewManagerFromURL(databaseURL, "./internal/migrations", false)
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
	log.Printf("Database connection test successful!")
	log.Printf("Health status: %t", status.IsHealthy)
	log.Printf("Pool stats - Total: %d, Idle: %d, Acquired: %d", 
		status.PoolStats.TotalConns, 
		status.PoolStats.IdleConns, 
		status.PoolStats.AcquiredConns)

	return nil
}

// SetupDatabase creates database and runs migrations
func SetupDatabase(databaseURL string) error {
	log.Printf("Setting up database...")
	
	// Create database manager with auto-migrate
	mgr, err := NewManagerFromURL(databaseURL, "./internal/migrations", true)
	if err != nil {
		return fmt.Errorf("failed to create database manager: %w", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		mgr.Shutdown(ctx)
	}()

	// Wait for database to be healthy
	if err := mgr.WaitForHealth(time.Second * 30); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	log.Printf("Database setup completed successfully!")
	return nil
}
