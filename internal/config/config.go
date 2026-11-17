// Package config загружает конфигурацию приложения.
package config

//nolint:depguard
import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config описывает конфигурацию приложения.
type Config struct {
	Port          string
	PostgresURL   string
	RedisURL      string
	RedisTTL      time.Duration
	CBRAPIURL     string
	RatesFilePath string
}

// Load загружает конфигурацию из переменных окружения.
func Load() (*Config, error) {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Переменная PORT не задана. Скопируйте .env.example в .env и запустите снова")
	}

	pgURL := os.Getenv("PG_URL")
	if pgURL == "" {
		log.Fatal("Переменная PG_URL не задана. Скопируйте .env.example в .env и запустите снова")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("Переменная REDIS_URL не задана. Скопируйте .env.example в .env и запустите снова")
	}

	redisTTLStr := os.Getenv("REDIS_TTL")
	if redisTTLStr == "" {
		log.Fatal("Переменная REDIS_TTL не задана. Скопируйте .env.example в .env и запустите снова")
	}

	redisTTL, err := time.ParseDuration(redisTTLStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат REDIS_TTL: %w", err)
	}

	cbrURL := os.Getenv("CBR_URL")
	if cbrURL == "" {
		log.Fatal("Переменная CBR_URL не задана. Скопируйте .env.example в .env и запустите снова")
	}

	ratesFilePath := os.Getenv("RATES_FILE")
	if ratesFilePath == "" {
		log.Fatal("Переменная RATES_FILE не задана. Скопируйте .env.example в .env и запустите снова")
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
