package db

import (
    "com.ityurika/go-redis-clone/internal/data"
)


// SetString sets a string value for the given key.
func (kv *DB) SetString(key string, val string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.Store[key] = data.NewStringValue(val)
}