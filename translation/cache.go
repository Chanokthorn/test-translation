package translation

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string) (any, bool, error)
	Set(ctx context.Context, key string, value any) error
}

type cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) Cache {
	return &cache{
		client: client,
	}
}

func getKey(key string) string {
	return fmt.Sprintf("translation:%s", key)
}

func (c *cache) Get(ctx context.Context, key string) (any, bool, error) {
	val, err := c.client.Get(ctx, getKey(key)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, false, nil // Key does not exist
		}
		return nil, false, fmt.Errorf("failed to get value from cache: %w", err)
	}

	return val, true, nil
}

func (c *cache) Set(ctx context.Context, key string, value any) error {
	err := c.client.Set(ctx, getKey(key), value, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set value in cache: %w", err)
	}

	return nil
}
