package config_test

import (
	"testing"

	"github.com/ArtemQa8/currency-converter/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("Успешная загрузка конфига", func(t *testing.T) {
		t.Setenv("PORT", "9090")
		t.Setenv("PG_URL", "postgres://test:test@localhost:5432/test")
		t.Setenv("REDIS_URL", "localhost:6379")
		t.Setenv("REDIS_TTL", "12h")
		t.Setenv("CBR_URL", "http://test.com")
		t.Setenv("RATES_FILE", "./test.json")

		cfg, err := config.Load()

		assert.NoError(t, err)
		assert.Equal(t, "9090", cfg.Port)
		assert.Equal(t, "postgres://test:test@localhost:5432/test", cfg.PostgresURL)
		assert.Equal(t, "12h0m0s", cfg.RedisTTL.String())
	})
}
