package keeper

import (
	"context"

	"github.com/go-redis/redis/v8"
)

const NotFoundError = "not_found"

type Keeper interface {
	Get(key string) (string, error)
	Set(key string, message string) error
	Clean(key string) error
}

func GetDummyKeeper() Keeper {
	return DummyKeeper{make(map[string]string)}
}

func GetRedisKeeper() Keeper {
	return RedisKeeper{*redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}), context.Background()}
}
