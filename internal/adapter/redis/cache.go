package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) SetSession(ctx context.Context, token string, userID int64, expiration time.Duration) error {
	return c.client.Set(ctx, fmt.Sprintf("session:%s", token), userID, expiration).Err()
}

func (c *Cache) GetUserIDByToken(ctx context.Context, token string) (int64, error) {
	val, err := c.client.Get(ctx, fmt.Sprintf("session:%s", token)).Int64()
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (c *Cache) AcquireLock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return c.client.SetNX(ctx, fmt.Sprintf("lock:%s", key), "1", expiration).Result()
}

func (c *Cache) ReleaseLock(ctx context.Context, key string) error {
	return c.client.Del(ctx, fmt.Sprintf("lock:%s", key)).Err()
}
