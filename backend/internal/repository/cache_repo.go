package repository

import (
	"context"
	"fmt"
	"pdf-parser/pkg/database"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheRepository struct{}

func NewCacheRepository() *CacheRepository {
	return &CacheRepository{}
}

func (r *CacheRepository) SetTaskStatus(ctx context.Context, taskID, status string) error {
	key := fmt.Sprintf("task:%s:status", taskID)
	return database.GetRedis().Set(ctx, key, status, 24*time.Hour).Err()
}

func (r *CacheRepository) GetTaskStatus(ctx context.Context, taskID string) (string, error) {
	key := fmt.Sprintf("task:%s:status", taskID)
	return database.GetRedis().Get(ctx, key).Result()
}

func (r *CacheRepository) SetTaskProgress(ctx context.Context, taskID string, progress int) error {
	key := fmt.Sprintf("task:%s:progress", taskID)
	return database.GetRedis().Set(ctx, key, progress, 24*time.Hour).Err()
}

func (r *CacheRepository) GetTaskProgress(ctx context.Context, taskID string) (int, error) {
	key := fmt.Sprintf("task:%s:progress", taskID)
	val, err := database.GetRedis().Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

func (r *CacheRepository) DeleteTaskCache(ctx context.Context, taskID string) error {
	keys := []string{
		fmt.Sprintf("task:%s:status", taskID),
		fmt.Sprintf("task:%s:progress", taskID),
	}
	return database.GetRedis().Del(ctx, keys...).Err()
}

func (r *CacheRepository) SetRateLimit(ctx context.Context, ip string, count int64) error {
	key := fmt.Sprintf("ratelimit:%s", ip)
	return database.GetRedis().Set(ctx, key, count, time.Minute).Err()
}

func (r *CacheRepository) GetRateLimit(ctx context.Context, ip string) (int64, error) {
	key := fmt.Sprintf("ratelimit:%s", ip)
	val, err := database.GetRedis().Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

func (r *CacheRepository) IncRateLimit(ctx context.Context, ip string) (int64, error) {
	key := fmt.Sprintf("ratelimit:%s", ip)
	return database.GetRedis().Incr(ctx, key).Result()
}
