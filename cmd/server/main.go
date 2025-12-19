package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"com.ityurika/go-redis-clone/internal/command"
	"com.ityurika/go-redis-clone/internal/protocol"
)

func main() {
	addr := ":6379"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening on", addr)

	// 创建信号通道
	sig := make(chan os.Signal, 1)

	// 捕获 SIGINT / SIGTERM
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// 处理连接
	go func() {
		for {
			conn, err := ln.Accept()
			command.LoadAOF("/media/shiokou/DevRepo/DevHub/Projects/2025-myapp/redis-golang/go-redis-server/persistence/aof.txt", conn)
			//每次创建链接的时候都会重新加载aof
			if err != nil {
				return
			}
			go handle(conn)
		}
	}()

	// 阻塞等待信号
	<-sig

	// 收到信号后优雅关闭
	log.Println("shutting down...")
	ln.Close()
}

func handle(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		arr, err := protocol.ParseArray(reader)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("error parsing request: %v", err)
			return
		}

		if len(arr) == 0 {
			// nothing to do
			continue
		}

		// Dispatch: arr[0] is command name, pass the full array as args
		command.HandleCommand(conn, arr[0], arr)
	}
}
