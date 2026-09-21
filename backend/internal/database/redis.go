package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/wishwall/wishwall/internal/config"
	"github.com/wishwall/wishwall/internal/constants"
)

// OpenRedis 建立 Redis 连接并做连通性探测。
func OpenRedis(cfg *config.Config, logger *slog.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis %s: %w", cfg.RedisAddr, err)
	}
	logger.Info(constants.LogRedisConnected, "addr", cfg.RedisAddr)
	return client, nil
}
