package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"ruaymanagement/backend/external/FacebookConnect/config/database"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database database.Config
	Facebook FacebookConfig
	ChatBot  ChatBotConfig
	Logging  LoggingConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port    int    `mapstructure:"port"`
	Host    string `mapstructure:"host"`
	Mode    string `mapstructure:"mode"`
	Timeout struct {
		Read  time.Duration `mapstructure:"read"`
		Write time.Duration `mapstructure:"write"`
		Idle  time.Duration `mapstructure:"idle"`
	} `mapstructure:"timeout"`
}

// FacebookConfig holds Facebook API configuration
type FacebookConfig struct {
	VerifyToken  string `mapstructure:"verify_token"`
	PageToken    string `mapstructure:"page_token"`
	AppSecret    string `mapstructure:"app_secret"`
	GraphAPIURL  string `mapstructure:"graph_api_url"`
	WebhookURL   string `mapstructure:"webhook_url"`
	RateLimit    struct {
		RequestsPerSecond int           `mapstructure:"requests_per_second"`
		BurstSize         int           `mapstructure:"burst_size"`
		RetryAttempts     int           `mapstructure:"retry_attempts"`
		RetryDelay        time.Duration `mapstructure:"retry_delay"`
	} `mapstructure:"rate_limit"`
}

// ChatBotConfig holds ChatBot service configuration
type ChatBotConfig struct {
	URL     string        `mapstructure:"url"`
	Timeout time.Duration `mapstructure:"timeout"`
	Retry   struct {
		Attempts int           `mapstructure:"attempts"`
		Delay    time.Duration `mapstructure:"delay"`
	} `mapstructure:"retry"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Set defaults
	setDefaults()

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Override with environment variables
	viper.AutomaticEnv()

	var config Config

	// Server configuration
	config.Server.Port = getEnvInt("PORT", 8085)
	config.Server.Host = getEnvString("HOST", "0.0.0.0")
	config.Server.Mode = getEnvString("GIN_MODE", "release")
	config.Server.Timeout.Read = getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second)
	config.Server.Timeout.Write = getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second)
	config.Server.Timeout.Idle = getEnvDuration("SERVER_IDLE_TIMEOUT", 120*time.Second)

	// Database configuration
	config.Database = database.Config{
		Host:            getEnvString("DB_HOST", "localhost"),
		Port:            getEnvInt("DB_PORT", 5432),
		Name:            getEnvString("DB_NAME", "chatbot_mvp"),
		User:            getEnvString("DB_USER", "chatbot_user"),
		Password:        getEnvString("DB_PASSWORD", "chatbot_pass123"),
		SSLMode:         getEnvString("DB_SSL_MODE", "disable"),
		MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 300*time.Second),
	}

	// Facebook configuration
	config.Facebook.VerifyToken = getEnvString("FB_VERIFY_TOKEN", "")
	config.Facebook.PageToken = getEnvString("FB_PAGE_TOKEN", "")
	config.Facebook.AppSecret = getEnvString("FB_APP_SECRET", "")
	config.Facebook.GraphAPIURL = getEnvString("FB_GRAPH_API_URL", "https://graph.facebook.com/v18.0")
	config.Facebook.WebhookURL = getEnvString("FB_WEBHOOK_URL", "")
	config.Facebook.RateLimit.RequestsPerSecond = getEnvInt("FB_RATE_LIMIT_RPS", 10)
	config.Facebook.RateLimit.BurstSize = getEnvInt("FB_RATE_LIMIT_BURST", 20)
	config.Facebook.RateLimit.RetryAttempts = getEnvInt("FB_RETRY_ATTEMPTS", 3)
	config.Facebook.RateLimit.RetryDelay = getEnvDuration("FB_RETRY_DELAY", 1*time.Second)

	// ChatBot configuration
	config.ChatBot.URL = getEnvString("CHATBOT_CORE_URL", "http://localhost:8086")
	config.ChatBot.Timeout = getEnvDuration("CHATBOT_TIMEOUT", 30*time.Second)
	config.ChatBot.Retry.Attempts = getEnvInt("CHATBOT_RETRY_ATTEMPTS", 3)
	config.ChatBot.Retry.Delay = getEnvDuration("CHATBOT_RETRY_DELAY", 1*time.Second)

	// Logging configuration
	config.Logging.Level = getEnvString("LOG_LEVEL", "info")
	config.Logging.Format = getEnvString("LOG_FORMAT", "json")

	return &config, nil
}

func setDefaults() {
	viper.SetDefault("server.port", 8085)
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.mode", "release")
	viper.SetDefault("facebook.graph_api_url", "https://graph.facebook.com/v18.0")
	viper.SetDefault("facebook.rate_limit.requests_per_second", 10)
	viper.SetDefault("facebook.rate_limit.burst_size", 20)
	viper.SetDefault("facebook.rate_limit.retry_attempts", 3)
	viper.SetDefault("facebook.rate_limit.retry_delay", "1s")
	viper.SetDefault("chatbot.timeout", "30s")
	viper.SetDefault("chatbot.retry.attempts", 3)
	viper.SetDefault("chatbot.retry.delay", "1s")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
}

func getEnvString(key, defaultValue string) string {
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

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	// Validate required configuration
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	return &cfg, nil
}

// validate checks if all required configuration is present
func (c *Config) validate() error {
	if c.VerifyToken == "" {
		return fmt.Errorf("FB_VERIFY_TOKEN is required")
	}
	
	if c.PageToken == "" {
		return fmt.Errorf("FB_PAGE_TOKEN is required")
	}
	
	if c.ChatBotCoreURL == "" {
		return fmt.Errorf("CHATBOT_CORE_URL is required")
	}
	
	return nil
}

// GetPort returns the port with colon prefix for server binding
func (c *Config) GetPort() string {
	if c.Port[0] != ':' {
		return ":" + c.Port
	}
	return c.Port
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.GinMode == "debug" || os.Getenv("APP_ENV") == "development"
}
