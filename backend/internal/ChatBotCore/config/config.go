package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config holds all configuration for the ChatBotCore service
type Config struct {
	// Server Configuration
	Port     string `mapstructure:"PORT"`
	GinMode  string `mapstructure:"GIN_MODE"`
	LogLevel string `mapstructure:"LOG_LEVEL"`

	// Database Configuration
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSL_MODE"`
	DBSchema   string `mapstructure:"DB_SCHEMA"`

	// External Services
	FacebookConnectURL string `mapstructure:"FACEBOOK_CONNECT_URL"`
	FacebookPageToken  string `mapstructure:"FACEBOOK_PAGE_TOKEN"`
	FacebookVerifyToken string `mapstructure:"FACEBOOK_VERIFY_TOKEN"`

	// Features
	EnableAnalytics bool `mapstructure:"ENABLE_ANALYTICS"`
	EnableLogging   bool `mapstructure:"ENABLE_LOGGING"`

	// Performance
	MaxConnections     int `mapstructure:"MAX_CONNECTIONS"`
	ConnectionTimeout  int `mapstructure:"CONNECTION_TIMEOUT"`
	RequestTimeout     int `mapstructure:"REQUEST_TIMEOUT"`
}

var globalConfig *Config

// Load initializes and loads configuration from environment variables and .env file
func Load() (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Read from environment variables
	v.AutomaticEnv()

	// Try to read from .env file (optional)
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("../../") // Go up to project root
	v.AddConfigPath("../../../") // Go up to project root from cmd/server

	if err := v.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %v", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Post-process configuration
	if err := postProcess(&config); err != nil {
		return nil, fmt.Errorf("failed to post-process config: %w", err)
	}

	globalConfig = &config
	return &config, nil
}

// MustLoad loads configuration and panics if it fails
func MustLoad() *Config {
	config, err := Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	return config
}

// Get returns the global configuration instance
func Get() *Config {
	if globalConfig == nil {
		return MustLoad()
	}
	return globalConfig
}

// setDefaults sets default values for configuration
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("PORT", "8086")
	v.SetDefault("GIN_MODE", "debug")
	v.SetDefault("LOG_LEVEL", "info")

	// Database defaults
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5433")
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "password")
	v.SetDefault("DB_NAME", "loyverse_cache")
	v.SetDefault("DB_SSL_MODE", "disable")
	v.SetDefault("DB_SCHEMA", "chatbot_mvp")

	// External services
	v.SetDefault("FACEBOOK_CONNECT_URL", "http://localhost:8085")

	// Features
	v.SetDefault("ENABLE_ANALYTICS", true)
	v.SetDefault("ENABLE_LOGGING", true)

	// Performance
	v.SetDefault("MAX_CONNECTIONS", 25)
	v.SetDefault("CONNECTION_TIMEOUT", 30)
	v.SetDefault("REQUEST_TIMEOUT", 30)
}

// postProcess handles any post-processing of configuration values
func postProcess(config *Config) error {
	// Ensure required environment variables are set
	if config.DBPassword == "" {
		if dbPass := os.Getenv("DB_PASSWORD"); dbPass != "" {
			config.DBPassword = dbPass
		}
	}

	return nil
}

// GetDatabaseURL returns the complete database connection URL
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBSSLMode,
		c.DBSchema,
	)
}

// GetFacebookConnectURL returns the Facebook Connect service URL
func (c *Config) GetFacebookConnectURL() string {
	return c.FacebookConnectURL
}

// GetFacebookPageToken returns the Facebook Page token
func (c *Config) GetFacebookPageToken() string {
	return c.FacebookPageToken
}

// GetFacebookVerifyToken returns the Facebook webhook verify token
func (c *Config) GetFacebookVerifyToken() string {
	return c.FacebookVerifyToken
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("PORT is required")
	}
	if c.DBHost == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.DBPort == "" {
		return fmt.Errorf("DB_PORT is required")
	}
	if c.DBUser == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.DBPassword == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if c.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.FacebookConnectURL == "" {
		return fmt.Errorf("FACEBOOK_CONNECT_URL is required")
	}
	return nil
}
