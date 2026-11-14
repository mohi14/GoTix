package config

import (
	"log"
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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		DATABASE_URL:    os.Getenv("DATABASE_URL"),
		PORT:            os.Getenv("PORT"),
		JWTSECRET:       os.Getenv("JWTSECRET"),
		GOOGLECLIENTID:  os.Getenv("GOOGLE_CLIENT_ID"),
		SUPERADMINEMAIL: os.Getenv("SUPER_ADMIN_EMAIL"),
	}
}
