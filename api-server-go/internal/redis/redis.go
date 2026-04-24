package redis

import (
	"context"
	"fmt"
	"time"

	"mochat-api-server/internal/config"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis(cfg config.RedisConfig) error {
	RDB = redis.NewClient(&redis.Options{
		Addr:            fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:        cfg.Auth,
		DB:              cfg.DB,
		MinIdleConns:    cfg.MinConns,
		PoolSize:        cfg.MaxConns,
		DialTimeout:     fmtDuration(cfg.ConnTimeout),
		ReadTimeout:     fmtDuration(cfg.ConnTimeout),
		WriteTimeout:    fmtDuration(cfg.ConnTimeout),
		PoolTimeout:     fmtDuration(cfg.WaitTimeout),
		ConnMaxIdleTime: fmtDuration(cfg.MaxIdleTime),
	})

	ctx := context.Background()
	if err := RDB.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect redis: %w", err)
	}

	return nil
}

func fmtDuration(seconds float64) time.Duration {
	return time.Duration(seconds * float64(time.Second))
}

func CloseRedis() error {
	if RDB != nil {
		return RDB.Close()
	}
	return nil
}
