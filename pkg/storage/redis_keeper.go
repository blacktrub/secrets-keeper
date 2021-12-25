package keeper

import (
	"context"
	"errors"
	"fmt"
	"time"

	"secrets-keeper/pkg/encrypt"

	"github.com/go-redis/redis/v8"
)

type RedisKeeper struct {
	cn  redis.Client
	ctx context.Context
}

func (k RedisKeeper) Get(key string) (string, error) {
	val, err := k.cn.GetDel(k.ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New(NotFoundError)
	}
	realMessage, err := encrypt.Decrypt(val)
	if err != nil {
		return "", err
	}
	return realMessage, err
}

func (k RedisKeeper) GetRaw(key string) (string, error) {
	val, err := k.cn.GetDel(k.ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New(NotFoundError)
	}
	return val, err
}

func (k RedisKeeper) Set(key string, message string, ttl int) error {
	encryptedMessage, err := encrypt.Encrypt(message)
	if err != nil {
		return err
	}
    seconds, err := time.ParseDuration(fmt.Sprintf("%ds", ttl))
	if err != nil {
		return err
	}
	return k.cn.Set(k.ctx, key, encryptedMessage, seconds).Err()
}
