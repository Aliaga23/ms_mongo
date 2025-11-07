package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
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
	// Cargar archivo .env si existe
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

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
