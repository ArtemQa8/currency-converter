package redis_test

//nolint:depguard
import (
	"context"
	"testing"
	"time"

	"github.com/ArtemQa8/currency-converter/internal/model"
	"github.com/ArtemQa8/currency-converter/internal/repository/redis"
	"github.com/alicebob/miniredis/v2"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRedisCache(t *testing.T) {
	s := miniredis.RunT(t)

	ctx := context.Background()
	cache := redis.NewCache(s.Addr(), 1*time.Minute)

	rate := model.Rate{
		BaseCurrency:  model.USD,
		QuoteCurrency: model.RUB,
		Value:         decimal.NewFromFloat(90.5),
		Date:          time.Now(),
		Source:        "cbr",
	}

	t.Run("Сохранение и чтение курса", func(t *testing.T) {

		err := cache.SaveRate(ctx, rate)
		assert.NoError(t, err)

		gotRate, err := cache.GetRate(ctx, model.USD, model.RUB)
		assert.NoError(t, err)

		assert.Equal(t, model.USD, gotRate.BaseCurrency)
		assert.Equal(t, "90.5", gotRate.Value.String())
		assert.Equal(t, "redis", gotRate.Source)
	})

	t.Run("Ошибка при чтении несуществующего курса", func(t *testing.T) {
		_, err := cache.GetRate(ctx, model.EUR, model.RUB)
		assert.Error(t, err)
	})
}
