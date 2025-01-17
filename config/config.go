package config

import (
	"errors"
	"os"
)

type Config struct {
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	JWTSecret string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		DBHost:    os.Getenv("DB_HOST"),
		DBPort:    os.Getenv("DB_PORT"),
		DBUser:    os.Getenv("DB_USER"),
		DBPass:    os.Getenv("DB_PASS"),
		DBName:    os.Getenv("DB_NAME"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	// Basic check for missing env
	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBUser == "" ||
		cfg.DBPass == "" || cfg.DBName == "" || cfg.JWTSecret == "" {
		return nil, errors.New("missing required environment variables")
	}

	return cfg, nil
}
