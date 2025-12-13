package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL      string
	TelegramBotToken string
	JWTSecret        string
	Port             string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://elite:elite@localhost:5432/elite?sslmode=disable"),
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		JWTSecret:        getEnv("JWT_SECRET", "default-secret-key-change-in-production"),
		Port:             getEnv("PORT", "8080"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

