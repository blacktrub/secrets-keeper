package keeper

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
)

const NotFoundError = "not_found"

type Keeper interface {
	Get(key string) (string, error)
	Set(key string, message string) error
}

func GetDummyKeeper() Keeper {
    return DummyKeeper{mem: make(map[string]string), mu: &sync.Mutex{}}
}

func GetRedisKeeper() Keeper {
	return RedisKeeper{*redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}), context.Background()}
}
