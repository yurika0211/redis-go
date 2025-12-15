package db

import (
	"com.ityurika/go-redis-clone/internal/data"
)

/**
 * SetHash sets a hash value for the given key.
 * Returns an empty string and a boolean indicating if the operation was successful.
 */
func (kv *DB) HSet(key string, val string) (string, bool) {
	// Deprecated: old signature, keep for compatibility but not used.
	return "", true
}

// Deprecated: previous HGet removed. Use HGetField(key, field) instead.

// HSetField sets a field->value in the hash stored at key.
// Returns true if set successfully.
func (kv *DB) HSetField(key string, field string, val string) bool {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	// if key not present, create new hash
	v, exists := kv.store[key]
	if !exists {
		h := data.NewHash()
		h.SetField(field, val)
		kv.store[key] = h
		return true
	}
	// if present and is hash, set field
	if h, ok := v.(*data.Hash); ok {
		h.SetField(field, val)
		return true
	}
	// existing key is not a hash
	return false
}

// HGetField gets a field value from the hash stored at key.
// Returns value and true if found.
func (kv *DB) HGetField(key string, field string) (string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	v, exists := kv.store[key]
	if !exists {
		return "", false
	}
	h, ok := v.(*data.Hash)
	if !ok {
		return "", false
	}
	val, ok := h.GetField(field)
	if !ok {
		return "", false
	}
	return val, true
}