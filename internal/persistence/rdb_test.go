package persistence

import (
	"testing"
	"com.ityurika/go-redis-clone/internal/db"
	"fmt"
)

func TestRDB(t *testing.T) {
	rdb := CreateRDBInstance("test.rdb")
	rdb.Save(db.GetDB())
	rdb.Load(db.GetDB())
	db.GetDB().GetStoreLength()
	fmt.Println("Test Terminated")
}

