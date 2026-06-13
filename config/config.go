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
// Expects DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE, and PORT to be set.
func Load() (*Config, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if host == "" || port == "" || user == "" || dbname == "" {
		return nil, fmt.Errorf("missing required database environment variables")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	appPort := os.Getenv("PORT")
	if appPort == "" {
		appPort = "8080"
	}

	return &Config{
		Port:  appPort,
		DBDSN: dsn,
	}, nil
}
