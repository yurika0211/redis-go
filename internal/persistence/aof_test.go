package persistence

import (
	"fmt"
	"testing"
)

func TestAOF(t *testing.T) {
	fmt.Println("Testing AOF")
	aof, _ := OpenAOF("test.aof")
	aof.Append([]string{"SET", "foo", "bar"})
	fmt.Println("Done")
}
