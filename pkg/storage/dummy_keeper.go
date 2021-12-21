package keeper

import (
	"errors"
	"secrets-keeper/pkg/encrypt"
	"sync"
)

type DummyKeeper struct {
	mem map[string]string
	mu  *sync.Mutex
}

func (k DummyKeeper) Get(key string) (string, error) {
	k.mu.Lock()
	defer k.mu.Unlock()
	value, ok := k.mem[key]
	if !ok {
		return "", errors.New(NotFoundError)
	}
	realMessage, err := encrypt.Decrypt(value)
	if err != nil {
		return "", err
	}

	k.Clean(key)
	return realMessage, nil
}

func (k DummyKeeper) GetRaw(key string) (string, error) {
	value, _ := k.mem[key]
	return value, nil
}

func (k DummyKeeper) Set(key string, message string, ttl int) error {
	encryptedMessage, err := encrypt.Encrypt(message)
	if err != nil {
		return err
	}

	k.mem[key] = encryptedMessage
	return nil
}

func (k DummyKeeper) Clean(key string) error {
	delete(k.mem, key)
	return nil
}
