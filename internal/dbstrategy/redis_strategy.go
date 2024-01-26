package dbstrategy

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStrategy struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisStrategy(address, password string, port int) *RedisStrategy {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", address, port),
		Password: password,
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
