package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ArtemQa8/currency-converter/internal/model"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCache(addr string, ttl time.Duration) *Cache {
	return &Cache{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
		ttl: ttl,
	}
}

func (c *Cache) SaveRate(ctx context.Context, rate model.Rate) error {
	key := fmt.Sprintf("%s:%s:rate", rate.BaseCurrency, rate.QuoteCurrency)

	data, err := json.Marshal(rate)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, c.ttl).Err()
}

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
