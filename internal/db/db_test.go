package db

import "testing"
func TestSetStringGetString(t *testing.T) {
	kv := GetDB()
	key := "mykey"
	value := "myvalue"

	kv.SetString(key, value)

	retrievedValue, exists := kv.GetString(key)
	if !exists {
		t.Fatalf("Expected key %q to exist", key)
	}
	if retrievedValue != value {
		t.Fatalf("Expected value %q, got %q", value, retrievedValue)
	}
}//实现对string数据类型的存储