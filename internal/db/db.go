package db

// |              internal/db/db.go             |
// |-------------------------------------------|
// |  - 定义 DB (map[string]Value)         |
// |  - 支持 String 类型的 Value                |
// |  - 提供接口：SetString(key, val)、GetString(key)

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

