package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the ChatBotCore service
type Config struct {
	Port      string
	Host      string
	DBHost    string
	DBPort    int
	DBName    string
	DBUser    string
	DBPass    string
	DBSSLMode string
	LogLevel  string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Port:      getEnv("PORT", "8090"),
		Host:      getEnv("HOST", "0.0.0.0"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnvInt("DB_PORT", 5432),
		DBName:    getEnv("DB_NAME", "chatbot_mvp"),
		DBUser:    getEnv("DB_USER", "chatbot_user"),
		DBPass:    getEnv("DB_PASSWORD", "chatbot_pass123"),
		DBSSLMode: getEnv("DB_SSL_MODE", "disable"),
		LogLevel:  getEnv("LOG_LEVEL", "info"),
	}

	// Validate required fields
	if cfg.DBName == "" {
		return nil, fmt.Errorf("DB_NAME is required")
	}
	if cfg.DBUser == "" {
		return nil, fmt.Errorf("DB_USER is required")
	}
	if cfg.DBPass == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}

	return cfg, nil
}

// GetDSN returns database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName, c.DBSSLMode)
}

// GetServerAddr returns server address
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
