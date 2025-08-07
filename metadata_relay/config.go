package main

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/caarlos0/env/v11"
)

// holds all application configuration
type Config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
	Port     string `env:"PORT" envDefault:"8000"`

	CacheHost string `env:"VALKEY_HOST" envDefault:"localhost"`
	CachePort int    `env:"VALKEY_PORT" envDefault:"6379"`
	CacheDB   int    `env:"VALKEY_DB" envDefault:"0"`

	TMDBAPIKey string `env:"TMDB_API_KEY"`
	TVDBAPIKey string `env:"TVDB_API_KEY"`
}

var AppConfig Config

// parses environment variables and loads configuration
func LoadConfig() error {
	if err := env.Parse(&AppConfig); err != nil {
		return err
	}
	return nil
}

// converts string log level to slog.Level
func (c *Config) GetLogLevel() slog.Level {
	verbosity := strings.ToLower(c.LogLevel)
	switch verbosity {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// returns the server address in host:port format
func (c *Config) GetServerAddr() string {
	return ":" + c.Port
}

// returns the cache address in host:port format
func (c *Config) GetCacheAddr() string {
	return fmt.Sprintf("%s:%d", c.CacheHost, c.CachePort)
}
