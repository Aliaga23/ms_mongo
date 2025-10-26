package config

import (
	"os"
	"time"
)

type (
	Config struct {
		AppName       string
		AppVersion    string
		HTTPPort      string
		DatabaseURI   string
		DatabaseName  string
		JWTSecret     string
		JWTExpiration time.Duration
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{
		AppName:       os.Getenv("APP_NAME"),
		AppVersion:    os.Getenv("APP_VERSION"),
		HTTPPort:      os.Getenv("HTTP_PORT"),
		DatabaseURI:   os.Getenv("DATABASE_URI"),
		DatabaseName:  os.Getenv("DATABASE_NAME"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiration: 24 * time.Hour,
	}

	return cfg, nil
}
