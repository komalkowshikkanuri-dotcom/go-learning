package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type URLCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewURLCache(addr string, ttl time.Duration) *URLCache {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &URLCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *URLCache) Get(ctx context.Context, shortCode string) (string, error) {
	return c.client.Get(ctx, shortCode).Result()
}

func (c *URLCache) Set(ctx context.Context, shortCode, originalURL string) error {
	return c.client.Set(ctx, shortCode, originalURL, c.ttl).Err()
}

func (c *URLCache) Close() error {
	return c.client.Close()
}
