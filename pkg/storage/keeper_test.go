package keeper

import "testing"
import "sync"
import "secrets-keeper/pkg/encrypt"

func TestDummyKeeperSet(t *testing.T) {
	keeper := DummyKeeper{mem: make(map[string]string), mu: &sync.Mutex{}}
	key := "foo"
	value := "bar"
	keeper.Set(key, value, 0)
	savedMsg, _ := keeper.Get(key)
	if savedMsg != value {
		t.Error("bad memory value")
	}
}

func TestDummyKeeperGet(t *testing.T) {
	keeper := DummyKeeper{mem: make(map[string]string), mu: &sync.Mutex{}}
	key := "foo"
	value := "bar"
	encryptedValue, _ := encrypt.Encrypt(value)
	keeper.mem[key] = encryptedValue
	value_from_get, _ := keeper.Get(key)
	if value_from_get != value {
		t.Error("bad value from get", value_from_get, value)
	}
}

func TestDummyKeeperClean(t *testing.T) {
	keeper := DummyKeeper{mem: make(map[string]string), mu: &sync.Mutex{}}
	key := "foo"
	value := "bar"
	keeper.mem[key] = value
	keeper.Clean(key)
	_, ok := keeper.mem[key]
	if ok {
		t.Error("clean does not work")
	}
}
