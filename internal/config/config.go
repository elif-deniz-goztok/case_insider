// Package config loads application configuration from environment variables.
package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration values.
type Config struct {
	Port  string
	DBDSN string
}

// Load reads configuration from environment variables.
// Accepts DATABASE_URL directly (injected by Railway and most PaaS platforms),
// or falls back to individual DB_HOST / DB_PORT / DB_USER / DB_PASSWORD / DB_NAME / DB_SSLMODE vars.
func Load() (*Config, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		sslmode := os.Getenv("DB_SSLMODE")

		if host == "" || port == "" || user == "" || dbname == "" {
			return nil, fmt.Errorf("missing required database environment variables (set DATABASE_URL or DB_HOST/DB_PORT/DB_USER/DB_NAME)")
		}

		dsn = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user, password, host, port, dbname, sslmode,
		)
	}

	appPort := os.Getenv("PORT")
	if appPort == "" {
		appPort = "8080"
	}

	return &Config{
		Port:  appPort,
		DBDSN: dsn,
	}, nil
}
