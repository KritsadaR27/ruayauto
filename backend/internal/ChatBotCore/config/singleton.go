package config

import (
	"sync"
)

var (
	instance *Config
	once     sync.Once
)

// GetSingleton returns the singleton configuration instance
func GetSingleton() *Config {
	once.Do(func() {
		instance = MustLoad()
	})
	return instance
}
