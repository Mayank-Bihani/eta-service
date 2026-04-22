package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	DBUrl    string
	RedisUrl string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		Port:     os.Getenv("PORT"),
		DBUrl:    os.Getenv("DB_URL"),
		RedisUrl: os.Getenv("REDIS_URL"),
	}
}