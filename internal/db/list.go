package db

import (
	"com.ityurika/go-redis-clone/internal/data"
)

/**
 * LPUSH adds a value to the list stored at key.
 * Returns an empty slice and true if added successfully.
 */
func (kv *DB) LPUSH(key string, val string) ([]string, bool) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	v, exists := kv.store[key]
	if !exists {
		l := data.NewList([]string{})
		kv.store[key] = l
		v = l
	}
	if l, ok := v.(*data.List); ok {
		l.Push(val)
		return l.Values(), true
	}
	return nil, false
}

/**
 * LGET returns all values of the list stored at key.
 * Returns values and true if found.
 */
func (kv *DB) LGET(key string) ([]string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	v, exists := kv.store[key]
	if !exists {
		return nil, false
	}
	if l, ok := v.(*data.List); ok {
		return l.Values(), true
	}
	return nil, false
}