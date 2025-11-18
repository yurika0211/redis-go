package db

import (
	"sync"

	"com.ityurika/go-redis-clone/internal/data"
)

type DB struct {
	store map[string]data.Value
	mu    sync.RWMutex
}

var db *DB
var once sync.Once

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

// GetString retrieves the string value for the given key.
// Returns the value and a boolean indicating if the key exists.
func (kv *DB) GetString(key string) (string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	value, exists := kv.store[key]
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

/**
 * SADD adds a member to the set stored at key.
 * Returns true if added successfully.
 */
func (kv *DB) SADD(key string, member string) bool {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	
	// if key not present, create new set
	v, exists := kv.store[key]
	if !exists {
		s := data.NewSet(member)
		kv.store[key] = s
		return true
	}
	// if present and is set, add member
	if s, ok := v.(*data.Set); ok {
		// Assuming Set has an AddMember method
		s.AddMember(member)
		return true
	}
	// existing key is not a set
	return false
}

/**
 * SMEMBERS returns all members of the set stored at key.
 * Returns members and true if found.
 */
func (kv *DB) SMEMBERS(key string) ([]string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	v, exists := kv.store[key]
	if !exists {
		return nil, false
	}
	s, ok := v.(*data.Set)
	if !ok {
		return nil, false
	}
	// Assuming Set has a Members method that returns []string
	return s.Members(), true
}