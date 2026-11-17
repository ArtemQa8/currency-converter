package file_test

//nolint:depguard
import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/ArtemQa8/currency-converter/internal/model"
	"github.com/ArtemQa8/currency-converter/internal/repository/file"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestFileRepository(t *testing.T) {
	ctx := context.Background()

	setupRepo := func(t *testing.T) (*file.Repository, string) {
		t.Helper()
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test_rates.json")
		return file.NewRepository(filePath), filePath
	}

	t.Run("Сохранение и чтение", func(t *testing.T) {
		repo, filePath := setupRepo(t)

		rate := model.Rate{
			BaseCurrency:  model.USD,
			QuoteCurrency: model.RUB,
			Value:         decimal.NewFromFloat(90.5),
			Source:        "cbr",
		}

		err := repo.SaveRate(ctx, rate)
		assert.NoError(t, err)

		_, err = os.Stat(filePath)
		assert.NoError(t, err, "Файл должен существовать на диске")

		savedRate, err := repo.GetRate(ctx, model.USD, model.RUB)
		assert.NoError(t, err)
		assert.Equal(t, model.USD, savedRate.BaseCurrency)
		assert.Equal(t, "90.5", savedRate.Value.String())
		assert.Equal(t, "file", savedRate.Source)
	})

	t.Run("Обновление существующей валюты", func(t *testing.T) {
		repo, _ := setupRepo(t)

		rateUSD := model.Rate{BaseCurrency: model.USD, QuoteCurrency: model.RUB, Value: decimal.NewFromInt(90)}
		rateEUR := model.Rate{BaseCurrency: model.EUR, QuoteCurrency: model.RUB, Value: decimal.NewFromInt(100)}

		_ = repo.SaveRate(ctx, rateUSD)
		_ = repo.SaveRate(ctx, rateEUR)

		newRateUSD := model.Rate{BaseCurrency: model.USD, QuoteCurrency: model.RUB, Value: decimal.NewFromInt(95)}
		_ = repo.SaveRate(ctx, newRateUSD)

		gotUSD, _ := repo.GetRate(ctx, model.USD, model.RUB)
		assert.Equal(t, "95", gotUSD.Value.String())

		gotEUR, err := repo.GetRate(ctx, model.EUR, model.RUB)
		assert.NoError(t, err)
		assert.Equal(t, "100", gotEUR.Value.String())
	})

	t.Run("Чтение из seed-файла (если основного нет)", func(t *testing.T) {
		repo, _ := setupRepo(t)

		rate, err := repo.GetRate(ctx, model.USD, model.RUB)
		assert.NoError(t, err)
		assert.Equal(t, model.USD, rate.BaseCurrency)
		assert.Equal(t, "seed", rate.Source)
	})
}
