package persistence

import (
	"fmt"
	"testing"
)

func TestAOF(t *testing.T) {
	fmt.Println("Testing AOF")
	aof, err := OpenAOF("./log")
	if err != nil {
		panic(err)
	}
	// 手动指定策略，解决 nil panic
    aof.strategy = &NoStrategy{} 

    // 现在执行这一行就不会崩了
    if err := aof.Append([]string{"SET", "foo", "bar"}); err != nil {
        panic(err)
    }
	
	fmt.Println("Done")
}
