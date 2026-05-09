package database

import (
	"context"
	"fmt"
	"pdf-parser/internal/config"
	"time"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func InitRedis(cfg *config.RedisConfig) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	return nil
}

func GetRedis() *redis.Client {
	return rdb
}

func CloseRedis() error {
	if rdb != nil {
		return rdb.Close()
	}
	return nil
}
