// Package redis реализует кэш на базе Redis.
package redis

//nolint:depguard
import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ArtemQa8/currency-converter/internal/model"
	redislib "github.com/redis/go-redis/v9"
)

// Cache описывает кэш Redis.
type Cache struct {
	client *redislib.Client
	ttl    time.Duration
}

// NewCache конструктор - создаёт новый кэш.
func NewCache(addr string, ttl time.Duration) *Cache {
	return &Cache{
		client: redislib.NewClient(&redislib.Options{
			Addr: addr,
		}),
		ttl: ttl,
	}
}

// SaveRate сохраняет курс в кэш.
func (c *Cache) SaveRate(ctx context.Context, rate model.Rate) error {
	key := fmt.Sprintf("%s:%s:rate", rate.BaseCurrency, rate.QuoteCurrency)

	data, err := json.Marshal(rate)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, c.ttl).Err()
}

// GetRate получает курс из кэша.
func (c *Cache) GetRate(ctx context.Context, baseCurrency, quoteCurrency string) (model.Rate, error) {
	key := fmt.Sprintf("%s:%s:rate", baseCurrency, quoteCurrency)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return model.Rate{}, err
	}

	var rate model.Rate
	err = json.Unmarshal([]byte(val), &rate)
	if err != nil {
		return model.Rate{}, err
	}
	rate.Source = "redis"
	return rate, nil
}
