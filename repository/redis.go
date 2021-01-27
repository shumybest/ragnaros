package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/shumybest/ragnaros/config"
	"sync"
	"time"
)


type RedisClient struct {
	Redis *redis.Client
}

var redisInstance *RedisClient
var redisOnce sync.Once

func GetRedisInstance() *RedisClient {
	redisOnce.Do(func() {
		redisInstance = &RedisClient{}
	})
	return redisInstance
}

func (r *RedisClient) InitConnection() {
	if config.GetConfigString("spring.redis.host") != "" {
		r.Redis = redis.NewClient(&redis.Options{
			Addr:     config.GetConfigString("spring.redis.host") + ":" + config.GetConfigString("spring.redis.port"),
			Password: "",
			DB:       0,
		})
	}
}

func (r *RedisClient) RedisGet(key string) (string, error) {
	var ctx = context.Background()
	if r.Redis != nil {
		cache, err := r.Redis.Get(ctx, key).Result()

		if err != nil && err != redis.Nil {
			logger.Error("Get " + key + " error: " + err.Error())
			return "", err
		}

		if err == redis.Nil {
			return "", redis.Nil
		}

		return cache, nil
	} else {
		return "", redis.Nil
	}
}

func (r *RedisClient) RedisSet(key string, value interface{}, expiration time.Duration) {
	var ctx = context.Background()
	if r.Redis != nil {
		r.Redis.Set(ctx, key, value, expiration)
	}
}
