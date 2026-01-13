package db

import (
	"sync"
	"com.ityurika/go-redis-clone/internal/data"
)

type DB struct {
	Store map[string]data.Value
	mu    sync.RWMutex
}

var db *DB
var once sync.Once


/*
**GetDB returns the singleton instance of DB.
*/
func GetDB() *DB {
	once.Do(func() {
		db = &DB{
			Store: make(map[string]data.Value),
		}
	})
	return db
}

// GetString retrieves the string value for the given key.
// Returns the value and a boolean indicating if the key exists.
func (kv *DB) GetString(key string) (string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	value, exists := kv.Store[key]
	if !exists {
		return "", false
	}
	strValue, ok := value.(*data.StringValue)
	if !ok {
		return "", false
	}
	return strValue.Val, true
}

/**
 * GetStoreLength returns the number of key-value pairs in the store.
 */

func (kv *DB) GetStoreLength() int {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	return len(kv.Store)
}



