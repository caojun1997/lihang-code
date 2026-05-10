package common

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/songquanpeng/go-api-starter/common/logger"
)

var RDB redis.Cmdable
var RedisEnabled = true

func InitRedisClient() (err error) {
	if os.Getenv("REDIS_CONN_STRING") == "" {
		RedisEnabled = false
		logger.SysLog("REDIS_CONN_STRING not set, Redis is not enabled")
		return nil
	}

	redisConnString := os.Getenv("REDIS_CONN_STRING")
	opt, err := redis.ParseURL(redisConnString)
	if err != nil {
		logger.Fatal("failed to parse Redis connection string: " + err.Error())
	}

	logger.SysLog("Redis is enabled")
	RDB = redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = RDB.Ping(ctx).Result()
	if err != nil {
		logger.Fatal("Redis ping test failed: " + err.Error())
	}

	return nil
}

func RedisSet(key string, value string, expiration time.Duration) error {
	if !RedisEnabled {
		return nil
	}
	ctx := context.Background()
	return RDB.Set(ctx, key, value, expiration).Err()
}

func RedisGet(key string) (string, error) {
	if !RedisEnabled {
		return "", nil
	}
	ctx := context.Background()
	return RDB.Get(ctx, key).Result()
}

func RedisDel(key string) error {
	if !RedisEnabled {
		return nil
	}
	ctx := context.Background()
	return RDB.Del(ctx, key).Err()
}

func RedisDecrease(key string, value int64) error {
	if !RedisEnabled {
		return nil
	}
	ctx := context.Background()
	return RDB.DecrBy(ctx, key, value).Err()
}

func ParseRedisOption() *redis.Options {
	opt, err := redis.ParseURL(os.Getenv("REDIS_CONN_STRING"))
	if err != nil {
		logger.Fatal("failed to parse Redis connection string: " + err.Error())
	}
	return opt
}

func _() {
	_ = strings.Split("", ",")
}
