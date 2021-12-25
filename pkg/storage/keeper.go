package keeper

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
)

const NotFoundError = "not_found"

type Keeper interface {
	Get(key string) (string, error)
	GetRaw(key string) (string, error)
	Set(key string, message string, ttl int) error
}

func getRedisHost() string {
    redisHost := os.Getenv("STORAGE_HOST")
    if redisHost == "" {
        redisHost = "localhost"
    }
    return redisHost
}

func GetDummyKeeper() Keeper {
	return DummyKeeper{mem: make(map[string]string), mu: &sync.Mutex{}}
}

func GetRedisKeeper() Keeper {
	return RedisKeeper{*redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", getRedisHost()),
		Password: "",
		DB:       0,
	}), context.Background()}
}
