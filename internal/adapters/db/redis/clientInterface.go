package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis"
)

type RedisClient interface {
	Set(context.Context, string, interface{}, time.Duration) *redis.StatusCmd
}
