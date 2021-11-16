package keeper

import (
	"errors"
	"sync"
)

type DummyKeeper struct {
	mem map[string]string
	mu *sync.Mutex
}

func (k DummyKeeper) Get(key string) (string, error) {
    k.mu.Lock()
    defer k.mu.Unlock()
	value, ok := k.mem[key]
	if !ok {
		return "", errors.New(NotFoundError)
	}
	k.Clean(key)
	return value, nil
}

func (k DummyKeeper) Set(key string, message string) error {
	k.mem[key] = message
	return nil
}

func (k DummyKeeper) Clean(key string) error {
	delete(k.mem, key)
	return nil
}
