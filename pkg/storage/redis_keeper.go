package keeper

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

const TTL = 0

type RedisKeeper struct {
	cn  redis.Client
	ctx context.Context
}

func (k RedisKeeper) Get(key string) (string, error) {
	val, err := k.cn.GetDel(k.ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New(NotFoundError)
	}
	return val, err
}

func (k RedisKeeper) Set(key string, message string, ttl int) error {
	return k.cn.Set(k.ctx, key, message, time.Duration(ttl)).Err()
}

