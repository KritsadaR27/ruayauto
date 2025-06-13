package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB represents the database connection pool
type DB struct {
	Pool *pgxpool.Pool
}

// New creates a new database connection
func New(databaseURL string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// Close closes the database connection
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// HealthCheck checks if the database is healthy
func (db *DB) HealthCheck(ctx context.Context) error {
	if db.Pool == nil {
		return fmt.Errorf("database pool is nil")
	}

	if err := db.Pool.Ping(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// Ping tests the database connection
func (db *DB) Ping(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

// Stats returns current pool statistics
func (db *DB) Stats() *pgxpool.Stat {
	return db.Pool.Stat()
}

// NewFromURL creates a new database connection from URL (compatibility)
func NewFromURL(databaseURL string) (*DB, error) {
	return New(databaseURL)
}

// Config holds database configuration
type Config struct {
	Host            string
	Port            string
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxConnections  int32
	MinConnections  int32
	MaxLifetime     time.Duration
	MaxIdleTime     time.Duration
	ConnectTimeout  time.Duration
	HealthCheckTime time.Duration
}

// DefaultConfig returns a default database configuration
func DefaultConfig() *Config {
	return &Config{
		Host:            "localhost",
		Port:            "5432",
		User:            "ruay",
		Password:        "ruay1234",
		Database:        "ruayAutoMsg",
		SSLMode:         "disable",
		MaxConnections:  25,
		MinConnections:  5,
		MaxLifetime:     time.Hour,
		MaxIdleTime:     time.Minute * 30,
		ConnectTimeout:  time.Second * 10,
		HealthCheckTime: time.Minute,
	}
}

// NewWithConfig creates a new database connection with configuration
func NewWithConfig(config *Config) (*DB, error) {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.User, config.Password, config.Host, config.Port, config.Database, config.SSLMode)
	return New(databaseURL)
}

// QueryRow executes a query that is expected to return at most one row
func (db *DB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return db.Pool.QueryRow(ctx, sql, args...)
}

// Exec executes a query without returning any rows
func (db *DB) Exec(ctx context.Context, sql string, args ...interface{}) error {
	_, err := db.Pool.Exec(ctx, sql, args...)
	return err
}
