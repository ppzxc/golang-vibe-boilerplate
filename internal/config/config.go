// Package config provides application configuration loaded from environment variables.
package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
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
// Returns an error if required environment variables are missing.
func Load() (*Config, error) {
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
			Host:     os.Getenv("DATABASE_HOST"),
			Port:     dbPort,
			User:     os.Getenv("DATABASE_USER"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			Name:     os.Getenv("DATABASE_NAME"),
			SSLMode:  getEnv("DATABASE_SSLMODE", "disable"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	slog.Info("config loaded",
		"server.port", cfg.Server.Port,
		"database.host", cfg.Database.Host,
		"database.name", cfg.Database.Name,
	)

	return cfg, nil
}

func (c *Config) validate() error {
	var missing []string
	if c.Database.Host == "" {
		missing = append(missing, "DATABASE_HOST")
	}
	if c.Database.User == "" {
		missing = append(missing, "DATABASE_USER")
	}
	if c.Database.Password == "" {
		missing = append(missing, "DATABASE_PASSWORD")
	}
	if c.Database.Name == "" {
		missing = append(missing, "DATABASE_NAME")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
