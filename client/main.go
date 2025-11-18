package main

import (
	"fmt"

	"com.ityurika/go-redis-client/core"
)

func main() {
	client := core.NewClient("127.0.0.1:6378")
	core.StartClient(client)
	defer core.CloseClient(client)

	buf := make([]byte, 1024)
	n, _ := client.Conn.Read(buf)
	fmt.Println("Received from server:", string(buf[:n]))
}
