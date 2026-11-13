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
}

func NewCache(addr string) *Cache {
	return &Cache{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (c *Cache) SaveRate(ctx context.Context, rate model.Rate) error {
	key := fmt.Sprintf("%s:%s:rate", rate.BaseCurrency, rate.QuoteCurrency)

	data, err := json.Marshal(rate)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, 24*time.Hour).Err()
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

	return rate, nil
}
