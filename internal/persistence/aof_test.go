package persistence

import (
	"fmt"
	"testing"
)

func TestAOF(t *testing.T) {
	fmt.Println("Testing AOF")
	aof, err := OpenAOF("/media/shiokou/DevRepo/DevHub/Projects/2025-myapp/redis-golang/go-redis-server/log/aof.txt")
	if err != nil {
		panic(err)
	}
	aof.Append([]string{"SET", "foo", "bar"})
	fmt.Println("Done")
}
