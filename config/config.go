package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	WeatherApiUrl string
	DatabaseUrl   string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		WeatherApiUrl: os.Getenv("WEATHER_API_URL"),
		DatabaseUrl:   os.Getenv("DATABASE_URL"),
	}
}
