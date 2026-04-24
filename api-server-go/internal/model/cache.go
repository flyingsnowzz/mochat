package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type ModelCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewModelCache(rdb *redis.Client, ttl time.Duration) *ModelCache {
	return &ModelCache{rdb: rdb, ttl: ttl}
}

func (c *ModelCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (c *ModelCache) Set(ctx context.Context, key string, val interface{}) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, data, c.ttl).Err()
}

func (c *ModelCache) Delete(ctx context.Context, keys ...string) error {
	return c.rdb.Del(ctx, keys...).Err()
}

func (c *ModelCache) Remember(ctx context.Context, key string, dest interface{}, fn func() (interface{}, error)) error {
	err := c.Get(ctx, key, dest)
	if err == nil {
		return nil
	}

	result, err := fn()
	if err != nil {
		return err
	}

	if err := c.Set(ctx, key, result); err != nil {
		return err
	}

	data, _ := json.Marshal(result)
	return json.Unmarshal(data, dest)
}

func ModelCacheKey(prefix, table string, id uint) string {
	return fmt.Sprintf("mc:%s:m:%s:%d", prefix, table, id)
}
