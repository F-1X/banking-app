package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Banking DB
}

type DB struct {
	DSN string `env:"BANKING_DSN" require:"true"`
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed load .env: %+v", err)
	}
	dsn := os.Getenv("BANKING_DSN")
	if dsn == "" {
		log.Fatalf("failed load BANKING_DSN: %+v", err)
	}
	return &Config{
		Banking: DB{
			DSN: dsn,
		},
	}
}
