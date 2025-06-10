package config

import (
	"log"
	"sync"
)

var (
	instance *Config
	once     sync.Once
)

// Get returns the singleton configuration instance
func Get() *Config {
	once.Do(func() {
		var err error
		instance, err = Load()
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}
	})
	return instance
}

// MustLoad loads configuration and panics if it fails
func MustLoad() *Config {
	config, err := Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	
	if err := config.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}
	
	return config
}
