package postgres_test

//nolint:depguard
import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ArtemQa8/currency-converter/internal/model"
	"github.com/ArtemQa8/currency-converter/internal/repository/postgres"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

const testSource = "test"

func TestPostgresRepository(t *testing.T) {
	pgURL := os.Getenv("PG_URL")
	if pgURL == "" {
		pgURL = "postgres://converter:converter@localhost:5432/converter?sslmode=disable" //nolint:gosec
	}

	ctx := context.Background()
	repo, err := postgres.NewRepository(ctx, pgURL)
	if err != nil {
		t.Skip("Postgres недоступен. Пропускаем тест: ", err)
	}

	from := "TS1"
	to := "TS2"

	t.Run("Сохранение и получение последнего курса", func(t *testing.T) {
		rate := model.Rate{
			BaseCurrency:  from,
			QuoteCurrency: to,
			Value:         decimal.NewFromFloat(1.5),
			Date:          time.Now(),
			Source:        testSource,
		}

		err := repo.SaveRate(ctx, rate)
		assert.NoError(t, err)

		gotRate, err := repo.GetRate(ctx, from, to)
		assert.NoError(t, err)
		assert.Equal(t, "1.5", gotRate.Value.String())
		assert.Equal(t, "postgres", gotRate.Source)
	})

	t.Run("Получение истории", func(t *testing.T) {
		_ = repo.SaveRate(ctx, model.Rate{
			BaseCurrency:  from,
			QuoteCurrency: to,
			Value:         decimal.NewFromInt(10),
			Date:          time.Now().Add(-48 * time.Hour),
			Source:        testSource,
		})
		_ = repo.SaveRate(ctx, model.Rate{
			BaseCurrency:  from,
			QuoteCurrency: to,
			Value:         decimal.NewFromInt(20),
			Date:          time.Now().Add(-24 * time.Hour),
			Source:        testSource,
		})

		history, err := repo.GetHistory(ctx, from, to, 3)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(history), 3)
		if len(history) >= 2 {
			assert.True(t, history[0].Date.Before(history[1].Date))
		}
	})
}
