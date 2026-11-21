package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DATABASE_URL    string
	PORT            string
	JWTSECRET       string
	GOOGLECLIENTID  string
	SUPERADMINEMAIL string
}

func LoadEnv() *Config {
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgresql://postgres:postgres@localhost:5432/gotix?sslmode=disable"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DATABASE_URL:    databaseURL,
		PORT:            port,
		JWTSECRET:       os.Getenv("JWTSECRET"),
		GOOGLECLIENTID:  os.Getenv("GOOGLE_CLIENT_ID"),
		SUPERADMINEMAIL: os.Getenv("SUPER_ADMIN_EMAIL"),
	}
}
