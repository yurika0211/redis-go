package db

import (
    "com.ityurika/go-redis-clone/internal/data"
)

// GetDB returns the singleton instance of DB.
func GetDB() *DB {
	once.Do(func() {
		db = &DB{
			store: make(map[string]data.Value),
		}
	})
	return db
}

// SetString sets a string value for the given key.
func (kv *DB) SetString(key string, val string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.store[key] = data.NewStringValue(val)
}