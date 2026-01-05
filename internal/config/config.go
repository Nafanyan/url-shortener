// Package config provides configuration loading and management.
package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config represents the application configuration.
type Config struct {
	Env         string `yaml:"env" end-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

// HTTPServer contains HTTP server configuration.
type HTTPServer struct {
	Address         string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout         time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"5s"`
	User            string        `yaml:"user" env-required:"true"`
	Password        string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

// MustLoad loads configuration from file and environment variables.
// It panics if configuration cannot be loaded.
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("error opening config file: %s", err)
	}

	var config Config
	err := cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &config
}
