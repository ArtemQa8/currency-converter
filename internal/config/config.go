package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	PostgresURL   string
	RedisURL      string
	RedisTTL      time.Duration
	CBRAPIURL     string
	RatesFilePath string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	pgURL := os.Getenv("PG_URL")
	if pgURL == "" {
		pgURL = "postgres://converter:converter@localhost:5432/converter?sslmode=disable"
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	redisTTLStr := os.Getenv("REDIS_TTL")
	if redisTTLStr == "" {
		redisTTLStr = "24h"
	}

	redisTTL, err := time.ParseDuration(redisTTLStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат REDIS_TTL: %w", err)
	}

	cbrURL := os.Getenv("CBR_URL")
	if cbrURL == "" {
		cbrURL = "https://www.cbr-xml-daily.ru/daily_json.js"
	}

	ratesFilePath := os.Getenv("RATES_FILE")
	if ratesFilePath == "" {
		ratesFilePath = "./data/rates.json"
	}

	return &Config{
		Port:          port,
		PostgresURL:   pgURL,
		RedisURL:      redisURL,
		RedisTTL:      redisTTL,
		CBRAPIURL:     cbrURL,
		RatesFilePath: ratesFilePath,
	}, nil

}
