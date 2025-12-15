package db

import (
	"com.ityurika/go-redis-clone/internal/data"
)

/**
 * SADD adds a member to the set Stored at key.
 * Returns true if added successfully.
 */
func (kv *DB) SADD(key string, member string) bool {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	// if key not present, create new set
	v, exists := kv.Store[key]
	if !exists {
		s := data.NewSet(member)
		kv.Store[key] = s
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
 * SMEMBERS returns all members of the set Stored at key.
 * Returns members and true if found.
 */
func (kv *DB) SMEMBERS(key string) ([]string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	v, exists := kv.Store[key]
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