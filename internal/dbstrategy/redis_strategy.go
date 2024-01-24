package dbstrategy

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStrategy struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisStrategy() *RedisStrategy {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()

	return &RedisStrategy{Client: client, Ctx: ctx}
}

func (rs *RedisStrategy) GetKey(field string) string {
	value, _ := rs.Client.Get(rs.Ctx, field).Result()
	return value
}

func (rs *RedisStrategy) SetNewKeyWithTimeLimit(field string, timeLimit int) error {
	return rs.Client.Set(rs.Ctx, field, 1, time.Duration(timeLimit)*time.Second).Err()
}

func (rs *RedisStrategy) IncrementExistingkey(field string, currentRequests int) error {
	return rs.Client.Set(rs.Ctx, field, currentRequests, redis.KeepTTL).Err()
}
