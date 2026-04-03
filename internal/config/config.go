// Package config provides application configuration loaded from environment variables.
package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Port string
}

// DatabaseConfig holds PostgreSQL connection configuration.
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// DSN returns the PostgreSQL connection string.
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	dbPort := 5432
	if p := os.Getenv("DATABASE_PORT"); p != "" {
		if n, err := strconv.Atoi(p); err == nil {
			dbPort = n
		}
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", ":8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DATABASE_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DATABASE_USER", "postgres"),
			Password: getEnv("DATABASE_PASSWORD", "postgres"),
			Name:     getEnv("DATABASE_NAME", "boilerplate"),
			SSLMode:  getEnv("DATABASE_SSLMODE", "disable"),
		},
	}

	slog.Info("config loaded",
		"server.port", cfg.Server.Port,
		"database.host", cfg.Database.Host,
		"database.name", cfg.Database.Name,
	)

	return cfg
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
