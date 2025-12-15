package db

import (
	"sync"
	"fmt"
	"com.ityurika/go-redis-clone/internal/data"
)

type DB struct {
	store map[string]data.Value
	mu    sync.RWMutex
}

var db *DB
var once sync.Once


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
 * ZADD adds a member to the sorted set stored at key.
 * Returns "OK" and true if added successfully.
 */
func (kv *DB) ZADD (key string, score string, member string) (bool) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	// if key not present, create new sorted set
	v, exists := kv.store[key]
	if !exists {
		var val map[string]int
		z := data.NewSortedSet(val)
		z.Add(score, member)
		kv.store[key] = z
		return true
	}
	// if present and is sorted set, add member
	if z, ok := v.(*data.SortedSet); ok {
		z.Add(score, member)
		return true
	}
	// existing key is not a sorted set
	return false
}

func (kv *DB) ZRANGE(key string, satrt string, stop string) ([]string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	v, exists := kv.store[key]
	if !exists {
		return nil, false
	}
	z, ok := v.(*data.SortedSet)
	if !ok {
		return nil, false
	}
	// 打印所有内容
	var members []string
	for member, score := range z.Val {
		members = append(members, fmt.Sprintf("%s: %d", member, score))
	}
	fmt.Println("Sorted Set Members:")
	for _, m := range members {
		fmt.Println(m)
	}
	return members, true
}