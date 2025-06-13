package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Manager manages the database connection, health monitoring, and migrations
type Manager struct {
	db            *DB
	healthMonitor *HealthMonitor
	config        *Config
	migrationConfig *MigrationConfig
	mu            sync.RWMutex
	isShuttingDown bool
}

// ManagerConfig contains configuration for the database manager
type ManagerConfig struct {
	Database   *Config
	Migration  *MigrationConfig
	HealthCheck struct {
		Interval time.Duration
		Timeout  time.Duration
	}
	AutoMigrate bool
}

// DefaultManagerConfig returns a default manager configuration
func DefaultManagerConfig() *ManagerConfig {
	return &ManagerConfig{
		Database:  DefaultConfig(),
		Migration: &MigrationConfig{
			MigrationsPath: "./migrations",
			DatabaseName:   "ruayAutoMsg",
		},
		HealthCheck: struct {
			Interval time.Duration
			Timeout  time.Duration
		}{
			Interval: time.Minute,
			Timeout:  time.Second * 30,
		},
		AutoMigrate: true,
	}
}

// NewManager creates a new database manager
func NewManager(config *ManagerConfig) (*Manager, error) {
	// Create database connection
	db, err := NewWithConfig(config.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	// Create health monitor
	healthMonitor := NewHealthMonitor(
		db,
		config.HealthCheck.Interval,
		config.HealthCheck.Timeout,
	)

	manager := &Manager{
		db:              db,
		healthMonitor:   healthMonitor,
		config:          config.Database,
		migrationConfig: config.Migration,
	}

	// Run migrations if auto-migrate is enabled
	if config.AutoMigrate {
		if err := manager.runMigrations(); err != nil {
			// Don't fail completely if migrations fail, just log the error
			log.Printf("Warning: Failed to run migrations: %v", err)
		}
	}

	// Start health monitoring
	healthMonitor.Start()

	log.Printf("Database manager initialized successfully")
	return manager, nil
}

// NewManagerFromURL creates a new database manager from a database URL
func NewManagerFromURL(databaseURL string, migrationPath string, autoMigrate bool) (*Manager, error) {
	// Create database connection
	db, err := NewFromURL(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	// Create health monitor with default settings
	healthMonitor := NewHealthMonitor(
		db,
		time.Minute,
		time.Second*30,
	)

	manager := &Manager{
		db:            db,
		healthMonitor: healthMonitor,
		migrationConfig: &MigrationConfig{
			MigrationsPath: migrationPath,
			DatabaseName:   "ruayAutoMsg", // Default database name
		},
	}

	// Run migrations if auto-migrate is enabled
	if autoMigrate {
		if err := manager.runMigrations(); err != nil {
			// Don't fail completely if migrations fail, just log the error
			log.Printf("Warning: Failed to run migrations: %v", err)
		}
	}

	// Start health monitoring
	healthMonitor.Start()

	log.Printf("Database manager initialized from URL")
	return manager, nil
}

// runMigrations runs database migrations
func (m *Manager) runMigrations() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	// Check if migrations directory exists
	if _, err := os.Stat(m.migrationConfig.MigrationsPath); os.IsNotExist(err) {
		log.Printf("Migrations directory '%s' does not exist, skipping migrations", m.migrationConfig.MigrationsPath)
		return nil
	}

	return m.db.RunMigrations(ctx, m.migrationConfig)
}

// GetDB returns the database connection
func (m *Manager) GetDB() *DB {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.db
}

// GetHealthStatus returns the current health status
func (m *Manager) GetHealthStatus() HealthStatus {
	return m.healthMonitor.GetStatus()
}

// IsHealthy returns whether the database is currently healthy
func (m *Manager) IsHealthy() bool {
	return m.healthMonitor.IsHealthy()
}

// WaitForHealth waits for the database to become healthy
func (m *Manager) WaitForHealth(timeout time.Duration) error {
	return m.healthMonitor.WaitForHealth(timeout)
}

// RunMigrations manually runs migrations
func (m *Manager) RunMigrations() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	return m.db.RunMigrations(ctx, m.migrationConfig)
}

// RollbackMigration rolls back the last migration
func (m *Manager) RollbackMigration() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	return m.db.RollbackMigration(ctx, m.migrationConfig)
}

// Shutdown gracefully shuts down the database manager
func (m *Manager) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	if m.isShuttingDown {
		m.mu.Unlock()
		return nil
	}
	m.isShuttingDown = true
	m.mu.Unlock()

	log.Printf("Shutting down database manager...")

	// Stop health monitoring
	m.healthMonitor.Stop()

	// Close database connection
	m.db.Close()

	log.Printf("Database manager shutdown completed")
	return nil
}

// SetupGracefulShutdown sets up graceful shutdown handling
func (m *Manager) SetupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Printf("Received shutdown signal")

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		if err := m.Shutdown(ctx); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}

		os.Exit(0)
	}()
}

// Ping tests the database connection
func (m *Manager) Ping(ctx context.Context) error {
	return m.db.Ping(ctx)
}
